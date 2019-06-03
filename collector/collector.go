package collector

import (
	"bufio"
	"bytes"
	"github.com/bverschueren/check_mk_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	namespace = "check_mk"
)

var (
	factories  = make(map[string]func() (Collector, error))
	collectors = make(map[string]Collector)
	command    = "check_mk_agent"
)

func registerCollector(collector string, factory func() (Collector, error)) {
	factories[collector] = factory
}

type Collector interface {
	Update(unstructuredStats *[]string, ch chan<- prometheus.Metric) error
}

type CheckMKCollector struct {
	target     config.Target
	collectors map[string]Collector
	Command    string
}

func NewMKCheckCollector(sshtarget config.Target) (CheckMKCollector, error) {

	for collector, factory := range factories {
		var err error
		collectors[collector], err = factory()
		if err != nil {
			log.Errorf("Unable to initialize factory for collector '%s'", collector)
		}
	}
	return CheckMKCollector{
		target:     sshtarget,
		collectors: collectors,
		Command:    command,
	}, nil
}

func (mc CheckMKCollector) connect() (*ssh.Session, ssh.Conn, error) {

	log.Debugf("Trying identity file '%s'", mc.target.IdentityFile)
	log.Debugf("Trying user '%s'", mc.target.User)

	key, err := ioutil.ReadFile(mc.target.IdentityFile)
	if err != nil {
		log.Errorf("unable to read private key: %v", err)
		return nil, nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Errorf("unable to parse private key: %v", err)
		return nil, nil, err
	}

	config := &ssh.ClientConfig{
		User: mc.target.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", mc.target.HostName+":"+strconv.Itoa(mc.target.Port), config)
	if err != nil {
		log.Errorf("unable to connect: %v", err)
		return nil, nil, err
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
		return nil, nil, err
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session, connection, nil
}

func (mc CheckMKCollector) collectRawStats() (*bytes.Buffer, error) {
	log.Debugf("Collecting stats from %s", mc.target.HostName)
	var stdoutBuf bytes.Buffer

	session, connection, err := mc.connect()
	if err != nil {
		log.Infof("Unable to collect stats from '%s': %s", mc.target.HostName, err)
		return nil, err
	}

	session.Stdout = &stdoutBuf
	err = session.Run(mc.Command)

	log.Trace("Raw check_mk stats: " + stdoutBuf.String())

	session.Close()
	connection.Close()
	return &stdoutBuf, nil
}

// TODO: allow overriding subsystems
func (mc CheckMKCollector) structureRawStats(raw *bytes.Buffer) *map[string]*[]string {
	scanner := bufio.NewScanner(strings.NewReader(raw.String()))
	log.Trace("Raw stats: ", raw.String())

	re := regexp.MustCompile("<<<([\\w_]+)>>>")
	structuredStats := make(map[string]*[]string)

	curStat := "unknown"
	for scanner.Scan() {
		in := scanner.Text()
		// TODO: filter only configured subsystems
		if match := re.FindStringSubmatch(in); match != nil {
			log.Debugf("Parsing subsystem: %s", match[1])
			curStat = match[1]
			if _, exists := structuredStats[curStat]; !exists {
				structuredStats[curStat] = new([]string)
			}
		} else {
			log.Tracef("Parsing %s", in)
			*structuredStats[curStat] = append(*structuredStats[curStat], in)
		}
	}
	return &structuredStats
}

func (mc CheckMKCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(mc.collectors))
	if rawStats, err := mc.collectRawStats(); err == nil {
		structuredRawStats := mc.structureRawStats(rawStats)
		for name, c := range mc.collectors {
			log.Debugf("Collecting from '%s'", name)
			if _, ok := (*structuredRawStats)[name]; !ok {
				log.Debugf("No raw stats found for '%s'", name)
				continue
			}
			go func(name string, c Collector) {
				c.Update((*structuredRawStats)[name], ch)
				wg.Done()
			}(name, c)
		}
	}
	wg.Wait()
}

func (mc CheckMKCollector) Describe(ch chan<- *prometheus.Desc) {
	// TODO: implement metrics about the scrape
}
