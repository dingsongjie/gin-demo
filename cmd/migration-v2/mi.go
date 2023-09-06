package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/rjeczalik/notify"
)

func main() {
	var listenDir string
	godotenv.Load(".env")
	flag.StringVar(&listenDir, "listen-dir", "./...", "directory of file notify to listen")
	c := make(chan notify.EventInfo)

	if err := notify.Watch(listenDir, c, notify.Remove, notify.Create,
		notify.Write,
		notify.Rename); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Block until an event is received.
	for ei := range c {
		log.Println("Got event:", ei)
	}

}
