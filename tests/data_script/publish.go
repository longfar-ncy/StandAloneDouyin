package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	util "douyin/test/testutil"
)

const (
	username = "videos-publisher"
	password = "123456"

	VideoDir = "/tmp/videos/"
	VideoNum = 39
)

var (
	userid int64
	token  string
	// token = ""

	VideoNames [VideoNum]string
)

func main() {
	ch := make(chan int)
	go GetVideoData(ch)
	go GetUserInfo(ch)
	for i := 0; i < 2; i++ {
		<-ch
	}
	fmt.Println("Prepare over")

	PublishVideos()
}

func PublishVideos() {
	for i := 0; i < VideoNum; i++ {
		DoFeed(int32(i), token, fmt.Sprintf("video-%d", i))
	}
}

func GetUserInfo(ch chan<- int) {
	userid_, token_, err := util.GetUseridAndToken(username, password)
	assert(err)
	userid = userid_
	_ = token_
	ch <- 1
}

func buildBody(idx int32, token, title string) (*bytes.Buffer, string) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(VideoNames[idx])
	assert(err)
	defer file.Close()
	part1, err := writer.CreateFormFile("data", filepath.Base(VideoNames[idx]))
	assert(err)
	_, err = io.Copy(part1, file)

	_ = writer.WriteField("token", token)
	_ = writer.WriteField("title", title)
	head := writer.FormDataContentType()
	err = writer.Close()
	assert(err)
	return payload, head
}

func DoFeed(idx int32, token, title string) {
	body, _ := buildBody(idx, token, title)
	resp, err := http.Post(util.CreateURL("/douyin/publish/action", nil), "", body)
	assert(err)
	if resp.StatusCode != 200 {
		panic("response code != 200")
	}
	respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
	assert(err)
	if respData.StatusCode != 0 {
		panic("status code != 0")
	}
}

func GetVideoData(ch chan<- int) {
	files, err := os.ReadDir(VideoDir)
	assert(err)

	for i, f := range files {
		VideoNames[i] = VideoDir + f.Name()
	}

	ch <- 1
	fmt.Println("Got Videos data")
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
