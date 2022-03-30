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

// 前景  背景  颜色
// 30	 40	  黑色
// 31	 41	  红色
// 32	 42	  绿色
// 33	 43	  黄色
// 34	 44	  蓝色
// 35	 45	  紫色
// 36	 46	  深绿
// 37	 47	  白色

// 打印带有前景色的文本
func ColorsPrintF(message string, fg uint8, isNewLine bool) {
	end := ""
	if isNewLine {
		end = "\n"
	}
	fmt.Printf("\x1b[%dm%s\x1b[0m%s", fg, message, end)
}

// 打印带有背景色和前景色的文本
func ColorsPrintBF(message string, bg, fg uint8, isNewLine bool) {
	end := ""
	if isNewLine {
		end = "\n"
	}
	fmt.Printf("\x1b[%d;%dm%s\x1b[0m%s", bg, fg, message, end)
}

func init() {
	app := &cli.App{
		Version: "1.0",
		Name:    "bili-go",
		Usage:   "命令行中下载 bilibili 视频",
		Action: func(c *cli.Context) error {
			ColorsPrintBF("📺 BiliBili 视频下载! ", 44, 33, true)
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
			if match, err := regexp.MatchString("[B|b][V|v][0-9a-zA-Z]{10}\\b", BV); err == nil && match {
				break
			} else {
				ColorsPrintF("BV号错误!", 31, true)
			}

		}

		reader := bufio.NewReader(os.Stdin)
		ColorsPrintF("? ", 32, false)
		ColorsPrintF("请输入视频BV号: ", 37, false)
		bv, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		BV = strings.TrimSpace(bv)
	}

	for {
		if SavePath != "" {
			if fileInfo, err := os.Stat(SavePath); err == nil && fileInfo.IsDir() {
				break
			} else {
				ColorsPrintF("路径错误!", 31, true)
			}
		}

		reader := bufio.NewReader(os.Stdin)
		ColorsPrintF("? ", 32, false)
		ColorsPrintF("请输入视频存储路径(如果为空, 默认为当前路径): ", 37, false)
		path, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		if path == "" || path == "\r\n" || path == "\n" {
			path, err = os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		SavePath = strings.TrimSpace(path)
	}
}
