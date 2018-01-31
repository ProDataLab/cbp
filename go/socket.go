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

	"github.com/go-mangos/mangos"
)

type Socket struct {
	id            Id
	url           string
	socketType    string
	transportType string
	mangosSocket  mangos.Socket
}

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

// AddSocket creates a new component socket and returns it.
func NewSocket(socket string, transport string, url string) (Socket, error) {
	if !socketTypeStringMap.containsKey(socket) {
		return nil, ErrNoSocket
	}
	var stype = socketTypeStringMap[socket]
	switch stype {
	case REQ:
		sock, err = req.NewSocket()
	case REP:
		sock, err = rep.NewSocket()
	case PUB:
		sock, err = pub.NewSocket()
	case SUB:
		sock, err = sub.NewSocket()
	case PUSH:
		sock, err = push.NewSocket()
	case PULL:
		sock, err = pull.NewSocket()
	default:
		sock, err = nil, ErrNoSocket
	}
	setTransportType(transport)
	setURL(url)
	return sock, err
}

// SetSubscriptionFilters is used to set topic filters in a pub/sub protocol
func SetSubscriptionFilters(socket mangos.Socket, topics []byte) error {
	if socket.GetProtocol().Name() != "sub" {
		return ErrWrongSocketType
	}
	return sock.SetOption(mangos.OptionSubscribe, topics)
}

func setTransportType(transport string) (Socket, error) {
	if !transportTypeStringMap.containsKey(transport) {
		return nil, ErrNoSocket
	}
	var ttype = transportTypeStringMap[transport]
	switch ttype {
	case INPROC:
		sock.AddTransport(inproc.NewTransport())
	case IPC:
		sock.AddTransport(ipc.NewTransport())
	case TCP:
		sock.AddTransport(tcp.NewTransport())
	case WS:
		sock.AddTransport(ws.NewTransport())
	default:
		sock, err = nil, ErrNoSocket
	}
	return sock, err
}

func setURL(url string) {
	urlString = url
}
