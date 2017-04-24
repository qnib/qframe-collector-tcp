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
		"collector.test.port": "10001",
		"collector.test.docker-host": "unix:///var/run/docker.sock",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	pde, err := qframe_collector_docker_events.New(qChan, *cfg, "test")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go pde.Run()
	time.Sleep(2*time.Second)
	p, err := qframe_collector_tcp.New(qChan, *cfg, "test")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2*time.Second)
	bg := qChan.Data.Join()
	for {
		qm := bg.Recv().(qtypes.QMsg)
		fmt.Printf("#### Received (remote:%s): %s\n", qm.Host, qm.Msg)
		//break

	}
}
