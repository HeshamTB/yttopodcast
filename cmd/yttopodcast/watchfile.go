package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

var watcherStop chan struct{}

func Watch(ctx context.Context, op fsnotify.Op, filename string, out chan struct{}, log *log.Logger) {
    var watcher *fsnotify.Watcher
    watcherStop = make(chan struct{})

    for {
        w, err := fsnotify.NewWatcher() 
        watcher = w
        if err != nil {
            e := errors.Join(errors.New("filewatch: "), err)
            log.Println(e.Error())
            time.Sleep(time.Second * 30)
            continue
        } 
        break
    }
    defer watcher.Close()

    for {
        err := watcher.Add(filename)
        if err == nil {
            break
        }
        e := errors.Join(errors.New("filewatch: "), err)
        log.Println(e.Error())
        time.Sleep(time.Second * 30)
    }

    
    log.Println("watching chanlist file")
    mainl:
    for {
        select {
        case event, ok := <- watcher.Events:
            log.Println("watchfile got error")
            if !ok {
                log.Println("watchfile events chan closed!")
                break mainl
            }
            if event.Has(fsnotify.Write) {
                out <- struct{}{}
            }

        case err, ok := <- watcher.Errors:
            log.Println("watchfile got error")
            log.Println(err.Error()) 
            if !ok {
                log.Println("watcher errors chan closed!")
                break mainl
            }

        case <- watcherStop:
            log.Println("watchfile got stop")
            break mainl
    }
    }

    log.Println("watchfile stopped")

}

func StopWatch() {
    log.Println("stopping watchfile")
    watcherStop <- struct{}{}
    close(watcherStop)
}
