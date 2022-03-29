package cmd

import (
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
			Usage:       "视频存储位置",
			Destination: &SavePath,
			Required:    false,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
