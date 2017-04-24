package main

import (
	"log"
	"fmt"
	"time"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-collector-tcp/lib"
	"github.com/qnib/qframe-collector-docker-events/lib"
)

func Run(qChan qtypes.QChan, cfg config.Config, name string) {
	p, _ := qframe_collector_tcp.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"collector.tcp.port": "10001",
		"collector.tcp.docker-host": "unix:///var/run/docker.sock",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	pde, err := qframe_collector_docker_events.New(qChan, *cfg, "docker-events")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go pde.Run()
	time.Sleep(2*time.Second)
	p, err := qframe_collector_tcp.New(qChan, *cfg, "tcp")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2*time.Second)
	bg := qChan.Data.Join()
	for {
		qm := bg.Recv().(qtypes.QMsg)
		if qm.Source == "tcp" {
			fmt.Printf("#### Received '%s' from enriched container: %v\n", qm.Msg, qm.Data)
			break
		}
	}
}
