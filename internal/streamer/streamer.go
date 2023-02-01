package streamer

import (
	"commutator/internal/config"
	"fmt"
	"net"
)

type Streamer struct {
	daemon net.Listener
	client net.Listener
	err    chan error
}

func (s *Streamer) Run() chan error {
	go s.acceptLoop(s.daemon, s.handleDaemonConn)
	go s.acceptLoop(s.client, s.handleClientConn)
	return s.err
}

func NewStreamer() (streamer *Streamer, err error) {
	dlis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Instance.InternalPort))
	if err != nil {
		return
	}
	clis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Instance.ExternalPort))
	if err != nil {
		return
	}
	streamer = &Streamer{
		daemon: dlis,
		client: clis,
		err:    make(chan error),
	}
	return
}

func (s *Streamer) acceptLoop(lis net.Listener, handler func(net.Conn)) {
	for {
		conn, errAccept := lis.Accept()
		if errAccept != nil {
			s.err <- errAccept
			return
		}
		go handler(conn)
	}
}
