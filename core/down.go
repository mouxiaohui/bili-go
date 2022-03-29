package core

import (
	"github.com/mouxiaohui/bili-go/cmd"
)

var (
	VIEW_URL     string = "https://api.bilibili.com/x/web-interface/view"
	PASSPORT_URL string = "https://passport.bilibili.com"
)

func Run() error {
	cmd.InitArguments()

	return nil
}
