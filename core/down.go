package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mouxiaohui/bili-go/cmd"

	"github.com/AlecAivazis/survey/v2"
)

var (
	BASE_URL string       = "https://api.bilibili.com/"
	CLIENT   *http.Client = &http.Client{Timeout: time.Duration(10) * time.Second}
)

func Run() error {
	cmd.InitArguments()

	videoInfo, err := getVideoInfo(cmd.BV)
	if err != nil {
		return err
	}
	if videoInfo.Aid == 0 {
		return errors.New("未找到视频!")
	}

	err = downloadVideo(&videoInfo)
	if err != nil {
		return err
	}

	return err
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
func downloadVideo(videoInfo *VideoInfo) error {
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

	videoIndex, audioIndex, ok := selectQuality(
		videoUrl.Dash.GetVideoQualitys(),
		videoUrl.Dash.GetAudioQualitys(),
	)

	return nil
}

// 选择视频，音频质量
func selectQuality(videoQualitys, audioQualitys []string) (videoIndex, audioIndex int, ok bool) {
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
				Message:  "选择音频质量: ",
				Options:  audioQualitys,
				VimMode:  true,
				PageSize: 10,
			},
		},
	}

	answers := struct {
		VideoQuality string
		AudioQuality string
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

// 根据value获取数组下标
func getArrayIndex[T string](arr *[]T, match T) (index int, ok bool) {
	for i, v := range *arr {
		if v == match {
			index = i
			ok = true
			return
		}
	}

	return
}
