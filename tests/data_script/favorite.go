package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	util "douyin/test/testutil"
)

const (
	password = "123456"
)

const (
	NumBigStar   = 10
	NumMidStar   = 100
	NumSmallStar = 300
	NumNormal    = 20000
	NumTotal     = NumBigStar + NumMidStar + NumSmallStar + NumNormal

	NameBigStar   = "big-star"
	NameMidStar   = "mid-star"
	NameSmallStar = "small-star"
	NameNormal    = "normal"

	NumBigFavorite = 2000
	NumMidFavorite   = 500
	NumSmallFavorite = 50
)

var (
	TotalUsers [NumTotal]UserInfo
	BigStars   []UserInfo
	MidStars   []UserInfo
	SmallStars []UserInfo
	Normals    []UserInfo
)

var (
	ch     = make(chan int)
	indexs [NumTotal]int
)

type UserInfo struct {
	Uid     int64
	Token   string
	VidList []int64
}

func init() {
	for i := 0; i < NumTotal; i++ {
		indexs[i] = i
	}
}

func main() {
	go GetNormalInfo()
	go GetSmallStarInfo()
	go GetMidStarInfo()
	go GetBigStarInfo()

	for i := 0; i < 4; i++ {
		<-ch
	}
	fmt.Println("Prepare over")

	go FavoriteForSmallStars()
	go FavoriteForMidStars()
	go FavoriteForBigStars()
	for i := 0; i < 3; i++ {
		<-ch
	}
}

func Sample(length, num int) []int {
	idxs := indexs[:length]
	rand.Shuffle(length, func(i, j int) {
		idxs[i], idxs[j] = idxs[j], idxs[i]
	})
	return idxs[:num]
}

func FavoriteForSmallStars() {
	for _, u := range SmallStars {
		for _, id := range u.VidList {
			n := rand.Intn(NumSmallFavorite/5) - NumSmallFavorite/10 + NumSmallFavorite
			idxs := Sample(NumNormal, n)
			for _, i := range idxs {
				DoFavorite(TotalUsers[i].Token, id)
			}
		}
	}

	ch <- 1
	fmt.Println("finish commenting for small stars")
}

func FavoriteForMidStars() {
	for _, u := range MidStars {
		for _, id := range u.VidList {
			n := rand.Intn(NumMidFavorite/5) - NumMidFavorite/10 + NumMidFavorite
			idxs := Sample(NumNormal+NumSmallStar, n)
			for _, i := range idxs {
				DoFavorite(TotalUsers[i].Token, id)
			}
		}
	}

	ch <- 1
	fmt.Println("finish commenting for mid stars")
}

func FavoriteForBigStars() {
	for _, u := range BigStars {
		for _, id := range u.VidList {
			n := rand.Intn(NumBigFavorite/5) - NumBigFavorite/10 + NumBigFavorite
			idxs := Sample(NumNormal+NumSmallStar+NumMidStar, n)
			for _, i := range idxs {
				DoFavorite(TotalUsers[i].Token, id)
			}
		}
	}

	ch <- 1
	fmt.Println("finish commenting for big stars")
}

func DoFavorite(token string, vid int64) {
	q := map[string]string{
		"token":        token,
		"video_id":     strconv.Itoa(int(vid)),
		"action_type":  "1",
	}
	resp, err := http.Post(util.CreateURL("/douyin/favorite/action/", q), "", nil)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("return code != 200")
	}

	respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
	assert(err)
	if respData.StatusCode != 0 {
		panic("status code != 0")
	}
}

func ShowUserProgress() {
	last := time.Now()
	idx := 0
	for idx < NumTotal {
		for TotalUsers[idx].Uid == 0 {
			time.Sleep(time.Second * 3)
		}
		idx++

		now := time.Now()
		if now.Sub(last) > time.Second*30 {
			fmt.Printf("Current progress: %d/%d\n", idx, NumTotal)
			last = now
		}
	}
}

func GetNormalInfo() {
	Normals = TotalUsers[:NumNormal]

	for i := 0; i < NumTotal; i++ {
		name := fmt.Sprintf("%s-%d", NameNormal, i)
		var err error
		Normals[i].Uid, Normals[i].Token, err = util.Login(name, password)
		assert(err)
	}

	ch <- 1
	fmt.Println("Got big stars info")
}

func GetVideoList(uid int64, token string) []int64 {
	q := map[string]string{
		"user_id": strconv.Itoa(int(uid)),
		"token":   token,
	}
	url := util.CreateURL("/douyin/publish/list/", q)
	resp, err := http.Get(url)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("return code != 200")
	}
	respData, err := util.GetDouyinResponse[util.DouyinPublishListResponse](resp)
	assert(err)
	if respData.StatusCode != 0 {
		panic("status code != 0")
	}
	vids := make([]int64, len(respData.VideoList))
	for i, v := range respData.VideoList {
		vids[i] = v.ID
	}
	return vids
}

func GetBigStarInfo() {
	BigStars = TotalUsers[NumTotal-NumBigStar:]
	if len(BigStars) != NumBigStar {
		panic("length wrong")
	}

	for i := 0; i < NumBigStar; i++ {
		name := NameBigStar + fmt.Sprintf("-%d", i)
		var err error
		BigStars[i].Uid, BigStars[i].Token, err = util.Login(name, password)
		assert(err)

		BigStars[i].VidList = GetVideoList(BigStars[i].Uid, BigStars[i].Token)
	}

	ch <- 1
	fmt.Println("Got big stars info")
}

func GetMidStarInfo() {
	MidStars = TotalUsers[NumNormal+NumSmallStar : NumTotal-NumBigStar]
	if len(MidStars) != NumMidStar {
		panic("length wrong")
	}

	for i := 0; i < NumMidStar; i++ {
		name := NameMidStar + fmt.Sprintf("-%d", i)
		var err error
		MidStars[i].Uid, MidStars[i].Token, err = util.Login(name, password)
		assert(err)

		MidStars[i].VidList = GetVideoList(MidStars[i].Uid, MidStars[i].Token)
	}

	ch <- 1
	fmt.Println("Got mid stars info")
}

func GetSmallStarInfo() {
	SmallStars = TotalUsers[NumNormal : NumTotal-NumBigStar-NumMidStar]
	if len(SmallStars) != NumSmallStar {
		panic("length wrong")
	}

	for i := 0; i < NumSmallStar; i++ {
		name := NameSmallStar + fmt.Sprintf("-%d", i)
		var err error
		SmallStars[i].Uid, SmallStars[i].Token, err = util.Login(name, password)
		assert(err)

		SmallStars[i].VidList = GetVideoList(SmallStars[i].Uid, SmallStars[i].Token)
	}

	ch <- 1
	fmt.Println("Got small stars info")
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
