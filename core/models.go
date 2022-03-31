package core

import (
	"fmt"
)

var (
	// 一般80以上的都要大会员，无法获取
	VIDEO_QUALITY map[int]string = map[int]string{
		6:   "240P 极速",
		16:  "360P 流畅",
		32:  "480P 清晰",
		64:  "720P 高清",
		74:  "720P 60帧",
		80:  "1080P 高清",
		112: "1080P 高码率",
		116: "1080P 60帧",
		120: "4k 超清",
		125: "HDR 真彩",
		127: "超高清 8K",
	}

	AUDIO_QUALITY map[int]string = map[int]string{
		30216: "64  kbps",
		30232: "128 kbps",
		30280: "320 kbps",
	}
)

type ResultInfo[D VideoUrl | VideoInfo] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    D      `json:"data"`
}

type VideoInfo struct {
	Bvid  string `json:"bvid"`
	Aid   int    `json:"aid"`
	Cid   int    `json:"cid"`
	Title string `json:"title"`
}

type VideoUrl struct {
	Dash Dash `json:"dash"`
}

type Dash struct {
	Duration int     `json:"duration"`
	Videos   []Video `json:"video"`
	Audios   []Audio `json:"audio"`
}

type Video struct {
	Id        int    `json:"id"`
	BaseUrl   string `json:"baseUrl"`
	BandWidth int    `json:"bandwidth"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Audio struct {
	Id      int    `json:"id"`
	BaseUrl string `json:"baseUrl"`
}

// 获取所有视频质量文字描述
func (d *Dash) GetVideoQualitys() []string {
	var qualitys []string
	for _, v := range d.Videos {
		qualitys = append(qualitys, v.getQuality())
	}

	return qualitys
}

// 获取所有音频质量的文字描述
func (d *Dash) GetAudioQualitys() []string {
	var qualitys []string
	for _, a := range d.Audios {
		qualitys = append(qualitys, a.getQuality())
	}

	return qualitys
}

// 获取视频质量的文字描述
func (v *Video) getQuality() string {
	if qua := VIDEO_QUALITY[v.Id]; qua == "" {
		return "UNKOWN"
	} else {
		return fmt.Sprintf("%s | %dx%d | BandWidth: %d", qua, v.Width, v.Height, v.BandWidth)
	}
}

// 获取音频质量的文字描述
func (a *Audio) getQuality() string {
	if qua := AUDIO_QUALITY[a.Id]; qua == "" {
		return "UNKOWN"
	} else {
		return qua
	}
}
