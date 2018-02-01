// TODO:
//	Look at the goroutines very carefully.. and FIX!
//  Look at everything.. lots of hackery here
//  Does a sub socket's option filters need to be set after "Dial()"???
//  mangos sockets need to be closed: socket.Close()

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
		sendChannel 	[]byte chan 
		recvChannel 	[]byte chan 

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

// NewSocket creates a new component socket and returns it.
func NewSocket(name string, socket string, transport string, url string) (*Socket, error) {
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
		sock.sendChannel = make([]byte)
		sock.recvChannel = make([]byte)
	case "rep":
		msock, err = rep.NewSocket()
		sock.sendChannel = make([]byte)
		sock.recvChannel = make([]byte)
	case "pub":
		msock, err = pub.NewSocket()
		sock.sendChannel = make([]byte)
	case "sub":
		msock, err = sub.NewSocket()
		sock.recvChannel = make([]byte)
	case "push":
		msock, err = push.NewSocket()
		sock.sendChannel = make([]byte)
	case "pull":
		msock, err = pull.NewSocket()
		sock.recvChannel = make([]byte)
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

func RunSocket(sockets []Socket) error {
	for socket := range sockets {
		var err error 
		switch socket.socketType {
		case "req":
			go runReq(socket)
		case "rep":
			go runRep(socket)
		case "pub":
			go runPub(socket)
		case "sub":
			go runSub(socket)
		case "push":
			go runPush(socket)
		case "pull":
			go runPull(socket)
		default:
			msock, err = nil, ErrWrongSocketType
		}
		return nil 
	}
	return ErrEmptySocketsArray
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

func isSocketType(socketType string) bool {
	for _, v := range socketTypes {
		if v == socketType {
			return true
		}
	}
	return false
}

func hasRecvChannel(sock Socket) {
	types := []string{
		"req"
		"rep",
		"sub",
		"pull"
	}
	for _, v := range types {
		if v == sock.socketType {
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

func runReq(sock Socket) error {
	var (
		err error
		msg []byte
	)
	if err = socket.mangosSocket.Dial(socket.url); err != nil {
		return err 
	}
	for {
		msg <-sock.sendChannel
		if err = sock.mangosSocket.Send(<-sock.sendChannel); err != nil {
			return err 
		}
		if msg, err = sock.mangosSocket.Recv(); err != nil {
			return err 
		}
		sock.recvChannel <-msg 
	}
	return nil 
}

func runRep(socket Socket) error {
	var (
		err error
		msg []byte
	)
	if err = socket.mangosSocket.Listen(socket.url); err != nil {
		return err 
	}
	for {
		if msg, err = socket.mangosSocket.Recv(); err != nil {
			return err 
		}
		socket.recvChannel <- msg 

		if err = socket.mangosSocket.Send(<-socket.sendChannel); err != nil {
			return err 
		}
	}
	return nil 
}

func runPub(socket Socket) error {
	var (
		err error
		msg []byte 
	)
	if err = socket.mangosSocket.Listen(socket.url); err != nil {
		return err 
	}
	for {
		if err = socket.mangosSocket.Send(<-socket.sendChannel); err != nil {
			return err 
		}
	}
	return nil 
}

func runSub(socket Socket) error {
	var (
		err error
		msg []byte
	)
	if err = socket.mangosSocket.Dial(socket.url); err != nil {
		return err
	}
	for {
		if msg, err = socket.mangosSocket.Recv(); err != nil {
			return err 
		}
		socket.recvChannel <- msg 
	}
	return nil 
}

func runPush(socket Socket) err {
	var (
		err error
		msg []byte 
	)
	if err = socket.mangosSocket.Dial(socket.url); err != nil {
		return err 
	}
	for {
		if err = socket.mangosSocket.Send(<-socket.sendChannel); err != nil {
			return err 
		}
	}
	return nil 
}

func runPull(socket Socket) error {
	var (
		err error
		msg []byte 
	)
	if err = socket.mangosSocket.Listen(socket.url); err != nil {
		return err 
	}
	for {
		if msg, err = socket.mangosSocket.Recv(); err != nil {
			return err
		}
		socket.recvChannel <- msg 
	}
	return nil 
}