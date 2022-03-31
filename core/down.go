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

	err = download(&videoInfo)
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

	if err = requestVideoUrl(&videoUrl.Dash.Videos[videoIndex].BaseUrl, videoInfo, fileFormat); err != nil {
		return err
	}
	requestAudioUrl(&videoUrl.Dash.Audios[audioIndex].BaseUrl, videoInfo)

	return nil
}

func requestVideoUrl(url *string, videoInfo *VideoInfo, fileFormat string) error {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return err
	}
	generateHeaders(req, &videoInfo.Bvid)

	resp, err := CLIENT.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := fmt.Sprintf("%s.%s", videoInfo.Title, fileFormat)
	file, err := os.Create(filepath.Join(cmd.SavePath, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

func requestAudioUrl(url *string, videoInfo *VideoInfo) {

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
				Options: []string{"mp4", "flv", "avi", "f4v", "wmv"},
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

func generateHeaders(req *http.Request, bvid *string) {
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
	req.Header.Add("origin", "https://www.bilibili.com")
	req.Header.Add("range", "bytes=0-")
	req.Header.Add("referer", fmt.Sprintf("https://www.bilibili.com/video/%s", *bvid))
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
