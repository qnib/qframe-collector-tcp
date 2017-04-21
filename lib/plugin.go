package qframe_collector_tcp

import (
	"bytes"
	"fmt"
	"os"
	"net"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
)

const (
	version = "0.1.0"
	pluginTyp = "collector"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, name, version),
	}
	return p, err
}

func (p *Plugin) Run() {
	host := p.CfgStringOr("bind-host", "127.0.0.1")
	port := p.CfgStringOr("bind-port", "11001")
	// Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		p.Log("error", fmt.Sprintln("Error listening:", err.Error()))
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	p.Log("info", fmt.Sprintln("Listening on " + host + ":" + port))
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

func (p *Plugin) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1048576)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		qm := qtypes.NewQMsg("tcp", p.Name)
		n := bytes.Index(buf, []byte{0})
		qm.Msg = string(buf[:n-1])
		qm.Host = conn.RemoteAddr().String()
		p.QChan.Data.Send(qm)
	}
	// Close the connection when you're done with it.
	conn.Close()
}