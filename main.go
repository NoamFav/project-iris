package main

import (
	"github.com/NoamFav/Iris/src/config"
	"github.com/NoamFav/Iris/src/loop"
	"log"
	"os"
)

func main() {

	err := os.MkdirAll(config.AudioDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create audio dir: %v", err)
	}
	loop.StartLoop()
}
