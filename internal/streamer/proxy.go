package streamer

import "net"

func proxy(from, to net.Conn) error {
	for {
		buf, errRead := readAndTrim(from)
		if errRead != nil {
			return errRead
		}
		{
			_, errWrite := to.Write(buf)
			if errWrite != nil {
				return errWrite
			}
		}
		if len(buf) == 1 && buf[0] == 0 {
			break
		}
	}
	return nil
}
