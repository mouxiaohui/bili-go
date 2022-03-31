package main

import (
	"log"

	"github.com/mouxiaohui/bili-go/core"
)

func main() {
	if err := core.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
