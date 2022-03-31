package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// 检查是否安装ffmpeg
func ffmpegVersion() error {
	cmd := exec.Command("ffmpeg", "-version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return errors.New("未找到ffmpeg, 请先安装!")
	}

	return nil
}

// 使用ffmpeg合并文件
func ffmpegMergeFile(fileList []*string, outFile string) error {
	arg := []string{}
	for _, fp := range fileList {
		arg = append(arg, "-i", *fp)
	}

	arg = append(arg, "-vcodec", "copy", "-acodec", "copy", outFile)
	cmd := exec.Command("ffmpeg", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", "文件合并失败", out.String()))
	}

	return nil
}
