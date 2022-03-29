package core

import (
	"fmt"

	"github.com/mouxiaohui/bili-go/cmd"
)

func Run() error {
	fmt.Println(cmd.BV)
	fmt.Println(cmd.SavePath)
	return nil
}
