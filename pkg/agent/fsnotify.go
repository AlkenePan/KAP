package agent

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type FileWatcherMsg struct {
	File  string
	Level string
	Msg   string
}

func GetNewFileWatcherMsgChan() chan FileWatcherMsg {
	return make(chan FileWatcherMsg)
}

func AddNewWatcher(path string, msgChan chan FileWatcherMsg) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	_ = watch.Add(path)

	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						msgChan <- FileWatcherMsg{
							File:  path,
							Level: "high",
							Msg:   "File Writed",
						}
					}
				}
			}
		}
	}()

	select {}
}
