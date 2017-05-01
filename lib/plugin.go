package qframe_collector_tcp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"net"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-inventory/lib"
)

const (
	version = "0.2.0"
	pluginTyp = "collector"
	dockerAPI = "v1.29"
)

type Plugin struct {
	qtypes.Plugin
	buffer chan interface{}
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
	for {
		select {
		case msg := <- p.buffer:
			switch msg.(type) {
			case IncommingMsg:
				im := msg.(IncommingMsg)
				qm := qtypes.NewQMsg("tcp", p.Name)
				qm.Msg = im.Msg
				req := qframe_inventory.NewIPContainerRequest(im.Host)
				p.QChan.Data.Send(req)
				cnt := <- req.Back
				qm.Host = strings.Trim(cnt.Name, "/")
				qm.Data = cnt
				p.QChan.Data.Send(qm)
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
			Host: addrTuple[0],
		}
		p.Log("info", fmt.Sprintf("Received TCP message '%s' from '%s'", im.Msg, im.Host))
		p.buffer <- im
	}
	// Close the connection when you're done with it.
	conn.Close()
}
