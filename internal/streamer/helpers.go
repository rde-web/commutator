package streamer

import (
	"commutator/internal/config"
	"errors"
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

func readAndTrim(conn net.Conn) ([]byte, error) {
	var buf []byte = make([]byte, config.Instance.BufferSize)
	n, errRead := conn.Read(buf)
	if n > 0 {
		buf = buf[:n]
	}
	return buf, errRead
}

func getDaemonDescription(conn net.Conn) (*description, error) {
	buf, errRead := readAndTrim(conn)
	if errRead != nil {
		return nil, errRead
	}
	var desc description
	{
		errUnmarshall := msgpack.Unmarshal(buf, &desc)
		if errUnmarshall != nil {
			return nil, errUnmarshall
		}
	}
	if len(desc.DaemonID) == 0 {
		return nil, errors.New("desc.DaemonID not given")
	}
	return &desc, nil
}
