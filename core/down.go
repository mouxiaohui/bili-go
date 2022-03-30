package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mouxiaohui/bili-go/cmd"
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
	url := fmt.Sprintf("%s%s?bvid=%s", BASE_URL, "x/web-interface/view", bvid)

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
			"%s%s?fnval=80&avid=%s&cid=%s",
			BASE_URL,
			"x/player/playurl",
			strconv.Itoa(videoInfo.Aid),
			strconv.Itoa(videoInfo.Cid),
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func readCloserToString(rc *io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(*rc)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
