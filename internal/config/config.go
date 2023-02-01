package config

var Instance config

func init() {
	Instance = config{
		SocketsPath:  ".",
		BufferSize:   256,
		InternalPort: 8080,
		ExternalPort: 8085,
		Debug:        true,
	}
}

type config struct {
	SocketsPath  string
	BufferSize   int
	InternalPort int
	ExternalPort int
	Debug        bool
}
