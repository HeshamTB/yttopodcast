package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

func main() {
    cmd := ytdlp.New().ExtractAudio().GetURL().GetDescription()
    res, err := cmd.Run(
        context.Background(), 
        "https://www.youtube.com/watch?v=JD9IQRlQyh0",
        "https://www.youtube.com/watch?v=ZwkNTwWJP5k",
        )

    if err != nil {
        fmt.Printf("err: %v\n", err)
        return
    }

    lines := strings.Split(res.Stdout, "\n")
    for _, line := range lines {
        fmt.Printf("line: %v\n", line)
        if strings.HasPrefix(line, "https://") {
            fmt.Println("probably link")
        }
    }

}
