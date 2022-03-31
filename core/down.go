package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mouxiaohui/bili-go/cmd"

	"github.com/AlecAivazis/survey/v2"
)

var (
	BASE_URL string = "https://api.bilibili.com/"
	BVID     string
	CLIENT   *http.Client = &http.Client{Timeout: time.Duration(10) * time.Second}
)

func Run() error {
	if err := ffmpegVersion(); err != nil {
		return err
	}

	cmd.InitArguments()

	videoInfo, err := getVideoInfo(cmd.BV)
	if err != nil {
		return err
	}
	if videoInfo.Aid == 0 {
		return errors.New("未找到视频!")
	}
	BVID = videoInfo.Bvid

	err = download(&videoInfo)
	if err != nil {
		return err
	}

	cmd.ColorsPrintF("下载完成!", 32, false)

	return nil
}

// 获取视频信息
func getVideoInfo(bvid string) (VideoInfo, error) {
	var resultInfo ResultInfo[VideoInfo]
	url := fmt.Sprintf("%sx/web-interface/view?bvid=%s", BASE_URL, bvid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resultInfo.Data, err
	}

	resp, err := CLIENT.Do(req)
	if err != nil {
		return resultInfo.Data, err
	}
	defer resp.Body.Close()

	bodyString, err := readCloserToString(&resp.Body)
	if err != nil {
		return resultInfo.Data, err
	}

	err = json.Unmarshal([]byte(bodyString), &resultInfo)
	if err != nil {
		return resultInfo.Data, err
	}

	return resultInfo.Data, nil
}

// 获取视频下载连接地址等信息
func getVideoUrl(url string) (VideoUrl, error) {
	var resultInfo ResultInfo[VideoUrl]
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resultInfo.Data, err
	}

	resp, err := CLIENT.Do(req)
	if err != nil {
		return resultInfo.Data, err
	}
	defer resp.Body.Close()

	bodyString, err := readCloserToString(&resp.Body)
	if err != nil {
		return resultInfo.Data, err
	}

	err = json.Unmarshal([]byte(bodyString), &resultInfo)
	if err != nil {
		return resultInfo.Data, err
	}

	return resultInfo.Data, nil
}

// 下载视频
func download(videoInfo *VideoInfo) error {
	videoUrl, err := getVideoUrl(
		fmt.Sprintf(
			"%sx/player/playurl?fnval=80&avid=%d&cid=%d",
			BASE_URL,
			videoInfo.Aid,
			videoInfo.Cid,
		),
	)
	if err != nil {
		return err
	}

	videoIndex, audioIndex, fileFormat, ok := selectQuality(
		videoUrl.Dash.GetVideoQualitys(),
		videoUrl.Dash.GetAudioQualitys(),
	)
	if !ok {
		return errors.New("未查询到选择项!")
	}

	audioFilePath := filepath.Join(
		cmd.SavePath,
		fmt.Sprintf("%d.audio", time.Now().Unix()),
	)
	videoFilePath := filepath.Join(
		cmd.SavePath,
		fmt.Sprintf("%d.video", time.Now().Unix()),
	)

	// 开启协程下载音频
	c := make(chan error)
	go func() {
		err = requestFileTo(&videoUrl.Dash.Audios[audioIndex].BaseUrl, audioFilePath)
		if err != nil {
			c <- err
		}

		c <- nil
	}()

	// 下载视频
	err = requestFileTo(&videoUrl.Dash.Videos[videoIndex].BaseUrl, videoFilePath)
	if err != nil {
		return err
	}

	outFile := filepath.Join(
		cmd.SavePath,
		fmt.Sprintf("%s_%s.%s", videoInfo.Title, getTimeFormat(), fileFormat),
	)

	// 等待协程
	err = <-c
	if err != nil {
		return err
	}

	// 合并视频/音频
	mergeFiles := []string{videoFilePath, audioFilePath}
	err = mergeAV(&outFile, &mergeFiles, &fileFormat)
	if err != nil {
		return err
	}

	// 删除合并前的文件
	if err := removeFiles(&mergeFiles); err != nil {
		return err
	}

	return nil
}

// 合并视频和音频
func mergeAV(outFile *string, mergeFiles *[]string, fileFormat *string) error {
	err := ffmpegMergeFile(
		mergeFiles,
		outFile,
	)
	if err != nil {
		// 如果合并失败，尝试合并成 MP4
		if *fileFormat != "mp4" {
			out := filepath.Join(
				cmd.SavePath,
				fmt.Sprintf("%s.%s", *outFile, "mp4"),
			)
			err = ffmpegMergeFile(
				mergeFiles,
				&out,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 请求数据, 并保存为指定文件格式
func requestFileTo(url *string, filePath string) error {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return err
	}
	generateHeaders(req)

	resp, err := CLIENT.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

// 选择视频，音频质量，视频保存格式
func selectQuality(videoQualitys, audioQualitys []string) (videoIndex, audioIndex int, fileFormat string, ok bool) {
	var qs = []*survey.Question{
		{
			Name: "VideoQuality",
			Prompt: &survey.Select{
				Message:  "选择视频画质: ",
				Options:  videoQualitys,
				VimMode:  true,
				PageSize: 10,
			},
		},
		{
			Name: "AudioQuality",
			Prompt: &survey.Select{
				Message: "选择音频质量: ",
				Options: audioQualitys,
				VimMode: true,
			},
		},
		{
			Name: "FileFormat",
			Prompt: &survey.Select{
				Message: "选择视频保存格式: ",
				Options: []string{"mp4", "avi", "wmv", "mov"},
				VimMode: true,
			},
		},
	}

	answers := struct {
		VideoQuality string
		AudioQuality string
		FileFormat   string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatal(err.Error())
	}

	vIndex, ok1 := getArrayIndex(&videoQualitys, answers.VideoQuality)
	aIndex, ok2 := getArrayIndex(&audioQualitys, answers.AudioQuality)
	if ok1 && ok2 {
		ok = true
		videoIndex = vIndex
		audioIndex = aIndex
		fileFormat = answers.FileFormat
	}

	return
}

func readCloserToString(rc *io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(*rc)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
