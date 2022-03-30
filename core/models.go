package core

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
	Id        int    `json:"id"`
	BaseUrl   string `json:"baseUrl"`
	BandWidth int    `json:"bandwidth"`
}

// 获取视频质量的文字描述
func GetQualityById(id int) string {
	// 一般80以上的都要大会员，无法获取
	switch id {
	case 6:
		return "240P 极速"
	case 16:
		return "360P 流畅"
	case 32:
		return "480P 清晰"
	case 64:
		return "720P 高清"
	case 74:
		return "720P60 高帧率"
	case 80:
		return "1080P 高清"
	case 112:
		return "1080P+ 高码率"
	case 116:
		return "1080P60 高帧率"
	case 120:
		return "4k 超清"
	case 125:
		return "HDR 真彩色"
	default:
		return "UNKNOW"
	}
}
