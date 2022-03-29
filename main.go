package main

import (
	"log"

	_ "github.com/mouxiaohui/bili-go/cmd"
	"github.com/mouxiaohui/bili-go/core"
)

func main() {
	if err := core.Run(); err != nil {
		log.Fatal(err)
	}
}
