package main

import (
	"fmt"
	"log"

	"github.com/gohugoio/hugo/watcher/filenotify"
)

func main() {
    
    w, err := filenotify.NewEventWatcher()    
    if err != nil {
        log.Fatalln(err.Error())
    }
    defer w.Close()

    err = w.Add("test")
    if err != nil {
        log.Fatalln(err.Error())
    }

    for {
        select {
        case e, ok:= <- w.Events():
            if !ok {
                log.Println("channel closed")
                return
            }
            fmt.Printf("e: %v\n", e)
    }
    }


}
