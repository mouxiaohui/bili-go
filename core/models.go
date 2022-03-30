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

type VideoUrl struct{}
