package main

import (
	"github.com/morawskim/asusrouter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func recordUnameMetrics(unameGauge *prometheus.GaugeVec) {
	go func() {
		for {
			uname, err := metric.UnameMetrics()
			if err != nil {
				log.Println(err)
			}

			unameGauge.With(prometheus.Labels{
				"domainname": uname.Domainname,
				"nodename":   uname.Nodename,
				"release":    uname.Release,
				"sysname":    uname.Sysname,
				"version":    uname.Version,
				"machine":    uname.Machine,
			}).Set(1)

			time.Sleep(1 * time.Second)
		}
	}()
}

func main() {
	bindTo := os.Args[1]

	unameGauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_uname_info",
		Help: "The node uname info",
	}, []string{
		"domainname",
		"nodename",
		"release",
		"sysname",
		"version",
		"machine",
	})

	recordUnameMetrics(unameGauge)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(bindTo+":2112", nil)
}
