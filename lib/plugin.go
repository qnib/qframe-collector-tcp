package qframe_collector_tcp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"net"
	"time"
	
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
)

const (
	version = "0.1.2"
	pluginTyp = "collector"
)

type Plugin struct {
	qtypes.Plugin
	buffer chan interface{}
	Inventory qtypes.ContainerInventory
}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, name, version),
		buffer: make(chan interface{}, 1000),
	}
	return p, err
}

func (p *Plugin) Run() {
	host := p.CfgStringOr("bind-host", "0.0.0.0")
	port := p.CfgStringOr("bind-port", "11001")
	p.Inventory = qtypes.NewContainerInventory()
	// Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		p.Log("error", fmt.Sprintln("Error listening:", err.Error()))
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	p.Log("info", fmt.Sprintln("Listening on " + host + ":" + port))
	go p.handleRequests(l)
	dc := p.QChan.Data.Join()
	//containers := map[string]types.ContainerJSON{}
	for {
		select {
		case msg := <- p.buffer:
			switch msg.(type) {
			case IncommingMsg:
				im := msg.(IncommingMsg)
				qm := qtypes.NewQMsg("tcp", p.Name)
				qm.Msg = im.Msg
				qm.Time = im.Time
				qm.Host = im.Host
				p.Log("debug", fmt.Sprintf("%s: %s", im.Host, qm.Msg))
				/*if cnt, ok := containers[im.Host]; ok {
					qm.Host = strings.Trim(cnt.Name, "/")
					qm.Data = cnt
					p.QChan.Data.Send(qm)
					continue
				}*/
				p.QChan.Data.Send(qm)

			}
		case dcMsg := <-dc.Read:
			switch dcMsg.(type) {
			case qtypes.QMsg:
				cm := dcMsg.(qtypes.QMsg)
				switch cm.Data.(type) {
				case qtypes.ContainerEvent:
					ce := cm.Data.(qtypes.ContainerEvent)
					p.Log("debug", fmt.Sprintf("%s.%s", ce.Event.Type, ce.Event.Action))
					/*for _, v := range ce.Container.NetworkSettings.Networks {
						containers[v.IPAddress] = ce.Container
					}*/

				}
			}
		}
	}

}

func (p *Plugin) handleRequests(l net.Listener) {
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go p.handleRequest(conn)
	}
}

type IncommingMsg struct {
	Msg string
	Time time.Time
	Host string
}

func (p *Plugin) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1048576)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		n := bytes.Index(buf, []byte{0})
		addrTuple := strings.Split(conn.RemoteAddr().String(), ":")
		im := IncommingMsg{
			Msg: string(buf[:n-1]),
			Time: time.Now(),
			Host: addrTuple[0],
		}
		p.buffer <- im
	}
	// Close the connection when you're done with it.
	conn.Close()
}
