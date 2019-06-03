package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bverschueren/check_mk_exporter/collector"
	"github.com/bverschueren/check_mk_exporter/config"
	"net/http"
	"strconv"
)

var (
	targets    = make(map[string]config.Target)
	cfg        = config.Config{
		Filename: kingpin.Flag(
			"config.file",
			"Config file to use",
		).Default("/etc/check_mk_exporter/ssh.yaml").String(),
	}
	listenPort = kingpin.Flag(
		"listen.port",
		"Port to listen on",
	).Default("2112").Int()
)

func checkMkHandler(w http.ResponseWriter, r *http.Request) {
	targetHost := r.URL.Query().Get("target")
	if targetHost == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}
	target, ok := targets[targetHost]
	if !ok {
		http.Error(w, fmt.Sprintf("Unknown target '%s'", targetHost), 400)
		log.Errorf("Unknown target '%s'", targetHost)
		return
	}

	// connection details overrides
	targetPortStr := r.URL.Query().Get("port")
	var targetPort int
	var err error
	if targetPortStr != "" {
		targetPort, err = strconv.Atoi(targetPortStr)
		if err == nil {
			target.Port = targetPort
		}
	}
	targetUser := r.URL.Query().Get("user")
	if targetUser != "" {
		target.User = targetUser
	}
	targetIdentityFile := r.URL.Query().Get("identityFile")
	if targetIdentityFile != "" {
		target.IdentityFile = targetIdentityFile
	}

	collector, _ := collector.NewMKCheckCollector(target)
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(w, r)
}

func init() {
	logLevel := kingpin.Flag(
		"log.level",
		"Enable specify log level",
	).Short('l').String()

	kingpin.Parse()

	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.Info("Enabling Debug mode.")
	case "trace":
		log.SetLevel(log.TraceLevel)
		log.Info("Enabling Trace mode.")
	}
}

func main() {
	kingpin.Parse()

	cfg.ReadFile(&targets)
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/check_mk", checkMkHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Check_MK Exporter</title></head>
             <body>
             <h1>Check_MK Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Infof("Start listening on :%d", *listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil))
}
