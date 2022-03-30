package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mouxiaohui/bili-go/cmd"
)

var (
	VIEW_URL     string       = "https://api.bilibili.com/x/web-interface/view"
	PASSPORT_URL string       = "https://passport.bilibili.com"
	CLIENT       *http.Client = &http.Client{Timeout: time.Duration(10) * time.Second}
)

func Run() error {
	cmd.InitArguments()

	videoInfo, err := getVideoInfo(cmd.BV)
	if err != nil {
		return err
	}
	if videoInfo.Aid == 0 {
		return errors.New("未找到视频❗")
	}

	return err
}

func getVideoInfo(bvid string) (VideoInfo, error) {
	var resultInfo ResultInfo[VideoInfo]
	url := fmt.Sprintf("%s?bvid=%s", VIEW_URL, bvid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resultInfo.Data, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := CLIENT.Do(req)
	if err != nil {
		return resultInfo.Data, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resultInfo.Data, err
	}
	bodyString := string(body)

	err = json.Unmarshal([]byte(bodyString), &resultInfo)
	return resultInfo.Data, err
}
