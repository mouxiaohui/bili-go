package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	BV       string
	SavePath string
)

func init() {
	app := &cli.App{
		Version: "1.0",
		Name:    "bili-go",
		Usage:   "命令行中下载 bilibili 视频",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "bv",
			Aliases:     []string{"b"},
			Usage:       "视频的bv号",
			Destination: &BV,
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "path",
			Aliases:     []string{"p"},
			Usage:       "视频存储位置(默认为当前路径)",
			Destination: &SavePath,
			Required:    false,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// 判断参数是否存在，如果不存在要求用户输入
func InitArguments() {
	if BV == "" {
		var bv string
		print("请输入视频BV号: ")
		fmt.Scan(&bv)
		BV = bv
	}

	for {
		if SavePath != "" {
			if fileInfo, err := os.Stat(SavePath); err == nil && fileInfo.IsDir() {
				break
			} else {
				println("路径不合法")
			}
		}

		var path string
		print("请输入视频存储路径: ")
		fmt.Scan(&path)
		SavePath = path
	}
}
