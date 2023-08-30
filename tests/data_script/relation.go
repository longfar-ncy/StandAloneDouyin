package main

import (
	"database/sql"
	util "douyin/test/testutil"
	"time"

	"fmt"
	"math/rand"
	"net/http"
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
)

var (
	TotalUsers [NumTotal]UserInfo
	BigStars   []UserInfo
	MidStars   []UserInfo
	SmallStars []UserInfo
	Normals    []UserInfo
	Mysql      *sql.DB
)

type UserInfo struct {
	Name  string
	Uid   int64
	Token string
}

func (f *UserInfo) Follow(u UserInfo) {
	query := map[string]string{
		"token":       f.Token,
		"action_type": "1",
	}
	query["to_user_id"] = fmt.Sprintf("%d", u.Uid)
	_, err := http.Post(util.CreateURL("/douyin/relation/action", query), "", nil)
	assert(err)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func ShowProgress() {
	idx := 0
	for {
		for TotalUsers[idx].Uid == 0 {
			time.Sleep(time.Second)
		}
		idx++

		if idx%100 == 0 {
			fmt.Printf("Current progress: %d/%d\n", idx, NumTotal)
		}
	}
}

func main() {
	var err error
	Mysql, err = util.GetDBConnection()
	assert(err)

	go ShowProgress()

	CreateNormal()
	CreateSmallStars()
	CreateMidStars()
	CreateBigStars()
}

func CreateNormal() {
	Normals = TotalUsers[:NumNormal]
	for i := 0; i < len(Normals); i++ {
		Normals[i].Name = fmt.Sprintf("%s-%d", NameNormal, i)
		var err error
		Normals[i].Uid, Normals[i].Token, err = util.GetUseridAndToken(Normals[i].Name, password)
		assert(err)
	}

	fmt.Println("create normals successfully!")
}

func CreateSmallStars() {
	SmallStars = TotalUsers[NumNormal : NumNormal+NumSmallStar]

	for i := 0; i < NumSmallStar; i++ {
		SmallStars[i].Name = fmt.Sprintf("%s-%d", NameSmallStar, i)
		var err error
		SmallStars[i].Uid, SmallStars[i].Token, err = util.GetUseridAndToken(SmallStars[i].Name, password)
		assert(err)

		if GetFansNum(Mysql, SmallStars[i]) >= 100 {
			continue
		}

		// add fans 100~1000
		n := rand.Intn(900) + 100
		fans := Sample(TotalUsers[:NumNormal], n)

		// do follow action
		for _, f := range fans {
			f.Follow(SmallStars[i])
		}
		fans = nil
	}

	fmt.Println("create small stars successfully!")
}

func CreateMidStars() {
	MidStars = TotalUsers[NumNormal+NumSmallStar : NumNormal+NumSmallStar+NumMidStar]

	for i := 0; i < NumMidStar; i++ {
		MidStars[i].Name = fmt.Sprintf("%s-%d", NameMidStar, i)
		var err error
		MidStars[i].Uid, MidStars[i].Token, err = util.GetUseridAndToken(MidStars[i].Name, password)
		assert(err)

		if GetFansNum(Mysql, MidStars[i]) >= 1000 {
			continue
		}

		// add fans 1000~7000
		n := rand.Intn(6000) + 1000
		fans := Sample(TotalUsers[:NumNormal+NumSmallStar], n)

		// do follow action
		for _, f := range fans {
			f.Follow(MidStars[i])
		}
		fans = nil
	}

	fmt.Println("create mid stars successfully!")
}

func CreateBigStars() {
	BigStars = TotalUsers[NumNormal+NumSmallStar+NumMidStar:]

	for i := 0; i < NumBigStar; i++ {
		BigStars[i].Name = fmt.Sprintf("%s-%d", NameBigStar, i)
		var err error
		BigStars[i].Uid, BigStars[i].Token, err = util.GetUseridAndToken(BigStars[i].Name, password)
		assert(err)

		if GetFansNum(Mysql, BigStars[i]) >= 10000 {
			continue
		}

		// add fans 10000~15000
		n := rand.Intn(5000) + 10000
		fans := Sample(TotalUsers[:NumTotal-NumBigStar], n)

		// do follow action
		for _, f := range fans {
			f.Follow(BigStars[i])
		}
		fans = nil
	}

	fmt.Println("create big stars successfully!")
}

func GetFansNum(db *sql.DB, user UserInfo) (res int) {
	// db, err := util.GetDBConnection()
	// assert(err)
	row, err := db.Query(`select follower_count from user where id=?`, user.Uid)
	assert(err)
	defer row.Close()
	row.Next()
	row.Scan(&res)
	return
}

func Sample(users []UserInfo, num int) []UserInfo {
	res := make([]UserInfo, len(users))
	copy(res, users)
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res[:num]
}

func CreateUsers(prefix string, num int) {
	for i := 1; i <= num; i++ {
		u := fmt.Sprintf("%s-%d", prefix, i)
		_, _, err := util.GetUseridAndToken(u, password)
		assert(err)
	}

	fmt.Printf("%d %s created", num, prefix)
}
