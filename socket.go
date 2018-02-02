// TODO:
//	Look at the goroutines very carefully.. and FIX!
//  Look at everything.. lots of hackery here
//  Does a sub socket's option filters need to be set after "Dial()"???
//  mangos sockets need to be closed: socket.Close()

package cbp

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
	"github.com/rs/xid"

	"github.com/go-mangos/mangos"
)

type (
	socket struct {
		id           _id
		url          string
		sockType     socketType
		transType    transportType
		mangosSocket mangos.Socket
		sendChannel  chan []byte
		recvChannel  chan []byte
	}
	socketType    string
	transportType string
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
	// ErrNoSocket blah
	ErrNoSocket = errors.New("socket creation failed")
	// ErrWrongSocketType blah
	ErrWrongSocketType = errors.New("wrong socket type")
	// ErrWrongTransportType blah
	ErrWrongTransportType = errors.New("wrong transport type")
)

// newSocket creates a new component socket and returns it.
func newSocket(name string, st socketType, tt transportType, url string) (*socket, error) {
	if !isSocketType(st) {
		return nil, ErrWrongSocketType
	}
	if !isTransportType(tt) {
		return nil, ErrWrongTransportType
	}
	s := new(socket)
	s.id.name = name
	s.id.uid = xid.New().String()
	s.url = url
	s.sockType = st
	s.transType = tt
	var (
		msock mangos.Socket
		err   error
	)
	switch st {
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
	s.mangosSocket = msock
	s.setTransportType(tt)
	return s, err
}

// setSubscriptionFilters is used to set topic filters in a pub/sub protocol
func (s *socket) setSubscriptionFilters(topics []byte) error {
	if s.mangosSocket.GetProtocol().Name() != "sub" {
		return ErrWrongSocketType
	}
	return s.mangosSocket.SetOption(mangos.OptionSubscribe, topics)
}

func runSocket(s socket) {
	switch s.sockType {
	case "req":
		s.sendChannel = make(chan []byte)
		s.recvChannel = make(chan []byte)
		go s.runReq()
	case "rep":
		s.sendChannel = make(chan []byte)
		s.recvChannel = make(chan []byte)
		go s.runRep()
	case "pub":
		s.sendChannel = make(chan []byte)
		go s.runPub()
	case "sub":
		s.recvChannel = make(chan []byte)
		go s.runSub()
	case "push":
		s.sendChannel = make(chan []byte)
		go s.runPush()
	case "pull":
		s.recvChannel = make(chan []byte)
		go s.runPull()
	}
}

func (s *socket) setTransportType(tt transportType) {
	switch tt {
	case "inproc":
		s.mangosSocket.AddTransport(inproc.NewTransport())
	case "ipc":
		s.mangosSocket.AddTransport(ipc.NewTransport())
	case "tcp":
		s.mangosSocket.AddTransport(tcp.NewTransport())
	case "ws":
		s.mangosSocket.AddTransport(ws.NewTransport())
	}
}

func isSocketType(st socketType) bool {
	for _, v := range socketTypes {
		if v == string(st) {
			return true
		}
	}
	return false
}

func (s *socket) hasRecvChannel() bool {
	recvTypes := []string{
		"req",
		"rep",
		"sub",
		"pull",
	}
	for _, v := range recvTypes {
		if v == string(s.sockType) {
			return true
		}
	}
	return false
}

func isTransportType(tt transportType) bool {
	for _, v := range transportTypes {
		if v == string(tt) {
			return true
		}
	}
	return false
}

func (s *socket) runReq() error {
	var (
		err error
		msg []byte
	)
	defer close(s.sendChannel)
	defer close(s.recvChannel)
	if err = s.mangosSocket.Dial(s.url); err != nil {
		return err
	}
	for {
		msg = <-s.sendChannel
		if err = s.mangosSocket.Send(<-s.sendChannel); err != nil {
			return err
		}
		if msg, err = s.mangosSocket.Recv(); err != nil {
			return err
		}
		s.recvChannel <- msg
	}
}

func (s *socket) runRep() error {
	var (
		err error
		msg []byte
	)
	defer close(s.recvChannel)
	defer close(s.sendChannel)
	if err = s.mangosSocket.Listen(s.url); err != nil {
		return err
	}
	for {
		if msg, err = s.mangosSocket.Recv(); err != nil {
			return err
		}
		s.recvChannel <- msg

		if err = s.mangosSocket.Send(<-s.sendChannel); err != nil {
			return err
		}
	}
}

func (s *socket) runPub() error {
	var err error
	defer close(s.sendChannel)
	if err = s.mangosSocket.Listen(s.url); err != nil {
		return err
	}
	for {
		if err = s.mangosSocket.Send(<-s.sendChannel); err != nil {
			return err
		}
	}
}

func (s *socket) runSub() error {
	var (
		err error
		msg []byte
	)
	defer close(s.recvChannel)
	if err = s.mangosSocket.Dial(s.url); err != nil {
		return err
	}
	for {
		if msg, err = s.mangosSocket.Recv(); err != nil {
			return err
		}
		s.recvChannel <- msg
	}
}

func (s *socket) runPush() error {
	var err error
	defer close(s.sendChannel)
	if err = s.mangosSocket.Dial(s.url); err != nil {
		return err
	}
	for {
		if err = s.mangosSocket.Send(<-s.sendChannel); err != nil {
			return err
		}
	}
}

func (s *socket) runPull() error {
	var (
		err error
		msg []byte
	)
	defer close(s.recvChannel)
	if err = s.mangosSocket.Listen(s.url); err != nil {
		return err
	}
	for {
		if msg, err = s.mangosSocket.Recv(); err != nil {
			return err
		}
		s.recvChannel <- msg
	}
}
