package cbp

import (
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/inproc"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/transport/ws"
	"github.com/rs/xid"

	"github.com/go-mangos/mangos"
)

type (
	Socket struct {
		id            Id
		url           string
		socketType    string
		transportType string
		mangosSocket  mangos.Socket
	}
)

var (
	socketTypes = [6]string{
		"req",
		"rep",
		"push",
		"pull",
		"pub",
		"sub",
	}
	transportTypes = [4]string{
		"inproc",
		"ipc",
		"tcp",
		"ws",
	}
)

func isSocketType(socketType string) bool {
	for _, v := range socketTypes {
		if v == socketType {
			return true
		}
	}
	return false
}

func isTransportType(transportType string) bool {
	for _, v := range transportTypes {
		if v == transportType {
			return true
		}
	}
	return false
}

// AddSocket creates a new component socket and returns it.
func NewSocket(name string, socket string, transport string, url string) (Socket, error) {
	if !isSocketType(socket) {
		return nil, ErrWrongSocketType
	}
	if !isSocketType(transport) {
		return nil, ErrWrongTransportType
	}
	sock := new(Socket)
	sock.id.name = name
	sock.id.uid = xid.New().String()
	sock.url = url
	sock.socketType = socket
	sock.transportType = transport
	var (
		msock mangos.Socket
		err error
	)
	switch socket {
	case "req":
		msock, err = req.NewSocket()
	case "rep":
		msock, err = rep.NewSocket()
	case "pub":
		msock, err = pub.NewSocket()
	case "sub":
		msock, err = sub.NewSocket()
	case "push":
		msock, err = push.NewSocket()
	case "pull":
		msock, err = pull.NewSocket()
	default:
		msock, err = nil, ErrNoSocket
	}
	if err != nil {
		return nil, err 
	}
	sock.mangosSocket = msock 
	setTransportType(sock, transport)
	return sock, err
}

// SetSubscriptionFilters is used to set topic filters in a pub/sub protocol
func SetSubscriptionFilters(socket mangos.Socket, topics []byte) error {
	if socket.GetProtocol().Name() != "sub" {
		return ErrWrongSocketType
	}
	return sock.SetOption(mangos.OptionSubscribe, topics)
}

func setTransportType(sock mangos.Socket, transport string) {
	switch transport {
	case "inproc":
		sock.AddTransport(inproc.NewTransport())
	case "ipc":
		sock.AddTransport(ipc.NewTransport())
	case "tcp":
		sock.AddTransport(tcp.NewTransport())
	case "ws":
		sock.AddTransport(ws.NewTransport())
	}
}

