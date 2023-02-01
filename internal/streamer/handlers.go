package streamer

import (
	"commutator/internal/config"
	"log"
	"net"
	"path"
)

func (s *Streamer) handleDaemonConn(conn net.Conn) {
	defer conn.Close()

	desc, errGetDD := getDaemonDescription(conn)
	if errGetDD != nil {
		if config.Instance.Debug {
			log.Println("handleDaemonConn get daemon desc", errGetDD)
		}
		return
	}
	var sockPath string = path.Join(config.Instance.SocketsPath, desc.DaemonID)
	lis, errListen := net.Listen("unix", sockPath)
	defer lis.Close()

	if errListen != nil {
		if config.Instance.Debug {
			log.Println("handleDaemonConn listen", errListen)
		}
		return
	}
	for {
		// @todo bottleneck. One request for daemon at once
		reqConn, errAccept := lis.Accept()
		if errAccept != nil {
			if config.Instance.Debug {
				log.Println("handleDaemonConn accept", errAccept)
			}
			continue
		}
		if err := proxy(reqConn, conn); err != nil {
			if config.Instance.Debug {
				log.Println("handleDaemonConn proxy1", err)
			}
			reqConn.Close()
			continue
		}
		if err := proxy(conn, reqConn); err != nil {
			if config.Instance.Debug {
				log.Println("handleDaemonConn proxy2", err)
			}
			reqConn.Close()
			continue
		}
	}
}

func (s *Streamer) handleClientConn(reqConn net.Conn) {
	defer reqConn.Close()
	desc, errGetDD := getDaemonDescription(reqConn)
	if errGetDD != nil {
		if config.Instance.Debug {
			log.Println("handleClientConn get daemon desc", errGetDD)
		}
		return
	}
	var sockPath string = path.Join(config.Instance.SocketsPath, desc.DaemonID)

	respConn, errDial := net.Dial("unix", sockPath)
	if errDial != nil {
		if config.Instance.Debug {
			log.Println("handleClientConn dial", errDial)
		}
		return
	}

	if err := proxy(reqConn, respConn); err != nil {
		if config.Instance.Debug {
			log.Println("handleClientConn proxy1", err)
		}
		return
	}
	if err := proxy(respConn, reqConn); err != nil {
		if config.Instance.Debug {
			log.Println("handleClientConn proxy2", err)
		}
		return
	}
}
