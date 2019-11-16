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

	for {
		select {
		case ev := <-watch.Events:
			{
				if ev.Op&fsnotify.Write == fsnotify.Write {
					msgChan <- FileWatcherMsg{
						File:  path,
						Level: "danger",
						Msg:   "File Writed",
					}
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					msgChan <- FileWatcherMsg{
						File:  path,
						Level: "danger",
						Msg:   "File Remove",
					}
				}
			}
		}
	}
}
