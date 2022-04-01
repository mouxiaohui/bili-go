package core

import (
	"fmt"
	"net/http"
	"os"
)

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

// 生成请求头文件
func generateHeaders(req *http.Request) {
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
	req.Header.Add("origin", "https://www.bilibili.com")
	req.Header.Add("range", "bytes=0-")
	req.Header.Add("referer", fmt.Sprintf("https://www.bilibili.com/video/%s", BVID))
}

// 删除文件
func removeFiles(fileList *[]string) error {
	for _, file := range *fileList {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
