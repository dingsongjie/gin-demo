package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dingsongjie/go-project-template/pkg/log"
	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	queue "github.com/nsqio/go-diskqueue"
	"github.com/rjeczalik/notify"
)

func main() {
	var listenDir string
	godotenv.Load(".env")
	flag.StringVar(&listenDir, "listen-dir", "./...", "directory of file notify to listen")
	log.Initialise()
	c := make(chan notify.EventInfo)

	if err := notify.Watch(listenDir, c, notify.Remove, notify.Create,
		notify.Write,
		notify.Rename); err != nil {
		log.Logger.Error(err.Error())
	}
	defer notify.Stop(c)

	mydir, _ := os.Getwd()
	dq := queue.New("notify-events", path.Join(mydir, "queue", "notify-events"), 1024, 4, 1<<10, 2500, 2*time.Second, NewLogger())
	defer dq.Close()
	// Block until an event is received.
	for ei := range c {
		log.Logger.Sugar().Infoln("Got event:", ei)
		buffer := bytes.Buffer{}
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(ei)
		if err != nil {
			log.Logger.Sugar().Error(ei.Event().String())
		}
		dq.Put(buffer.Bytes())
	}

}

func NewLogger() queue.AppLogFunc {
	return func(lvl queue.LogLevel, f string, args ...interface{}) {
		info := fmt.Sprintf(f, args...)
		if lvl == queue.DEBUG {
			log.Logger.Debug(info)
		} else if lvl == queue.ERROR {
			log.Logger.Error(info)
		} else if lvl == queue.FATAL {
			log.Logger.Fatal(info)
		} else if lvl == queue.INFO {
			log.Logger.Info(info)
		} else {
			log.Logger.Warn(info)
		}
	}
}
