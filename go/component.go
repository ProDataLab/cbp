package component

import (
	"errors"

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

type (
	// Socket is a mangos.Socket
	Socket         mangos.Socket
	socketsType    [6]string
	transportsType [4]string
)

var (
	// SocketTypes contains the socket type strings accepted
	SocketTypes = socketsType{
		"req",
		"rep",
		"push",
		"pull",
		"pub",
		"sub",
	}
	// TransportTypes contains the transport type strings accepted
	TransportTypes = transportsType{
		"inproc",
		"ipc",
		"tcp",
		"ws",
	}
	// ErrNoSocket is used when creating nanomsg socket fails
	ErrNoSocket = errors.New("can't aquire socket")
	// ErrWrongSocketType is used when setting filters on a push socket from a pull socket
	ErrWrongSocketType = errors.New("must be a pull socket")
	// ErrWrongTransportType is used when setting transport type with incorrect transport string
	ErrWrongTransportType = errors.New("can't set transport")
	urlString string
)

// NewSocket creates a new component socket and returns it.
func NewSocket(socket string, transport string, url string) (Socket, error) {
	if !SocketTypes.contains(socket) {
		return nil, ErrWrongSocketType
	} 
	if !TransportTypes.contains(transport) {
		return nil, ErrWrongTransportType
	}
	var (
		sock Socket
		err  error
	)
	switch socket {
	case "req":
		sock, err = req.NewSocket()
	case "rep":
		sock, err = rep.NewSocket()
	case "pub":
		sock, err = pub.NewSocket()
	case "sub":
		sock, err = sub.NewSocket()
	case "push":
		sock, err = push.NewSocket()
	case "pull":
		sock, err = pull.NewSocket()
	default:
		sock, err = nil, ErrNoSocket
	}
	setTransportType(sock, transport)
	setURL(url)
	return sock, err
}

// SetSubscriptionFilters is used to set topic filters in a pub/sub protocol
func SetSubscriptionFilters(sock Socket, topics []byte) error {
	if sock.GetProtocol().Name() != "sub" {
		return ErrWrongSocketType
	}
	return sock.SetOption(mangos.OptionSubscribe, topics)
}

func (st socketsType) contains(s string) bool {
	for _, v := range st {
		if s == v {
			return true
		}
	}
	return false
}

func (tt transportsType) contains(s string) bool {
	for _, v := range tt {
		if s == v {
			return true
		}
	}
	return false
}

func setTransportType(sock Socket, transport string) (Socket, error) {
	if !TransportTypes.contains(transport) {
		return nil, ErrWrongTransportType
	}
	var err error
	switch transport {
	case "inproc":
		sock.AddTransport(inproc.NewTransport())
	case "ipc":
		sock.AddTransport(ipc.NewTransport())
	case "tcp":
		sock.AddTransport(tcp.NewTransport())
	case "ws":
		sock.AddTransport(ws.NewTransport())
	default:
		sock, err = nil, ErrWrongTransportType
	}
	return sock, err
}

func setURL(url string) {
	urlString = url
}
