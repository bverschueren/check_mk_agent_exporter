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
func structureRawStats(raw *bytes.Buffer) *map[string]*[]string {
	scanner := bufio.NewScanner(strings.NewReader(raw.String()))
	log.Trace("Raw stats: ", raw.String())

	re := regexp.MustCompile("<<<([\\w_]+)>>>")
	structuredStats := make(map[string]*[]string)
	// list stats in temporary map to ensure unique elemets
	tempStats := make(map[string]struct{})
	curStat, prevStat := "", ""

	keyMapToList := func() *[]string {
		// convert map to slice of keys
		keys := new([]string)
		log.Tracef("Found for %s:", curStat)
		for k, _ := range tempStats {
			log.Trace(k)
			*keys = append(*keys, k)
		}
		return keys
	}

	for scanner.Scan() {

		// TODO: filter only configured subsystems
		in := scanner.Text()

		if match := re.FindStringSubmatch(in); match != nil {
			prevStat = curStat
			curStat = match[1]
			log.Debugf("Found stat %s", curStat)
			if prevStat == "" {
				prevStat = curStat
			}
			if prevStat != curStat {
				// move the stat from previous iteration to list
				structuredStats[prevStat] = keyMapToList()
				log.Debugf("structured %s", prevStat)
				tempStats = make(map[string]struct{})
			}
		} else {
			tempStats[in] = struct{}{}
		}
	}
	// process stat from the last iteration
	structuredStats[curStat] = keyMapToList()

	return &structuredStats
}

func (mc CheckMKCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(mc.collectors))
	if rawStats, err := mc.collectRawStats(); err == nil {
		structuredRawStats := structureRawStats(rawStats)
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
