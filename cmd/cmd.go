package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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
	for {
		if BV != "" {
			if match, err := regexp.MatchString("[B|b][V|v][0-9a-zA-Z]{10}", BV); err == nil && match {
				break
			} else {
				fmt.Println("BV号不合法")
			}

		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入视频BV号: ")
		bv, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		BV = strings.TrimSpace(bv)
	}

	for {
		if SavePath != "" {
			if fileInfo, err := os.Stat(SavePath); err == nil && fileInfo.IsDir() {
				break
			} else {
				fmt.Println("路径不合法")
			}
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入视频存储路径(如果为空, 默认为当前路径): ")
		path, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if path == "" || path == "\r\n" || path == "\n" {
			path, err = os.Getwd()
			fmt.Println("path: ", path)
			if err != nil {
				log.Fatal(err)
			}
		}
		SavePath = strings.TrimSpace(path)
	}
}
