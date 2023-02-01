package main

import (
	"commutator/internal/streamer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s, errMakeStreamer := streamer.NewStreamer()

	if errMakeStreamer != nil {
		panic(errMakeStreamer)
	}

	var sigch chan os.Signal = make(chan os.Signal, 3)

	signal.Notify(sigch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	select {
	case <-sigch:
		log.Println("exiting")
		return
	case e := <-s.Run():
		log.Println("error", e)
		return
	}
}
