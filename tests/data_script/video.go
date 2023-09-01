package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	util "douyin/test/testutil"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	TimeFormat    = "2023-09-01 12:00:00.000"
	FavorateCount = 0
	CommentCount  = 0
	NumUrls       = 44
)

type Video struct {
	Id            int    `gorm:"id;primarykey" json:"id"`                        // 自增主键
	PublishTime   string `gorm:"publish_time" json:"publish_time"`               // 视频发布时间
	AuthorId      int    `gorm:"author_id" json:"author_id"`                     // 作者用户ID
	PlayUrl       string `gorm:"play_url" json:"play_url"`                       // 视频播放URL
	CoverUrl      string `gorm:"cover_url" json:"cover_url"`                     // 封面图片URL
	FavoriteCount int    `gorm:"favorite_count;default:0" json:"favorite_count"` // 视频点赞数
	CommentCount  int    `gorm:"comment_count;default:0" json:"comment_count"`   // 视频评论数
	Title         string `gorm:"title" json:"title"`                             // 视频标题
}

type User struct {
	Id              int    `gorm:"id;primarykey" json:"id"`                             // 自增主键
	Username        string `gorm:"username" json:"username"`                 // 用户名，也是用户的昵称
	Password        string `gorm:"password" json:"password"`                 // 用户密码
	FollowingCount  int    `gorm:"following_count" json:"following_count"`   // 关注数
	FollowerCount   int    `gorm:"follower_count" json:"follower_count"`     // 粉丝数
	Avatar          string `gorm:"avatar" json:"avatar"`                     // 头像URL
	BackgroundImage string `gorm:"background_image" json:"background_image"` // 用户个人页顶部大图
	Signature       string `gorm:"signature" json:"signature"`               // 个人简介
	TotalFavorited  int    `gorm:"total_favorited" json:"total_favorited"`   // 获赞数量
	WorkCount       int    `gorm:"work_count" json:"work_count"`             // 作品数量
	FavoriteCount   int    `gorm:"favorite_count" json:"favorite_count"`     // 点赞视频数量
}

func (u User) TableName() string {
	return "user"
}

func (v Video) TableName() string {
	return "video"
}

const (
	password = "123456"
)

const (
	NumBigStar   = 10
	NumMidStar   = 100
	NumSmallStar = 300
	NumNormal    = 0
	NumTotal     = NumBigStar + NumMidStar + NumSmallStar + NumNormal

	NameBigStar   = "big-star"
	NameMidStar   = "mid-star"
	NameSmallStar = "small-star"
	NameNormal    = "normal"

	VideoNum = 39
	VideoDir = "/tmp/videos/"

	NumBigVideo   = 10
	NumMidVideo   = 5
	NumSmallVideo = 2

	NameBigVideo   = "Big Star's Video"
	NameMidVideo   = "Mid Star's Video"
	NameSmallVideo = "Small Star's Video"
)

var (
	PlayUrls  [NumUrls]string
	CoverUrls [NumUrls]string
)

var (
	db *gorm.DB
	ch = make(chan int)
)

var (
	TotalUsers [NumTotal]UserInfo
	BigStars   []UserInfo
	MidStars   []UserInfo
	SmallStars []UserInfo
	Normals    []UserInfo
)

type UserInfo struct {
	Uid   int64
	Token string
}

func init() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/douyin"))
	assert(err)
}

func main() {
	go GetSmallStarInfo()
	go GetMidStarInfo()
	go GetBigStarInfo()
	go GetUrls()

	for i := 0; i < 4; i++ {
		<-ch
	}

	go BigStarsPublish()
	go MidStarsPublish()
	go SmallStarsPublish()

	for i := 0; i < 3; i++ {
		<-ch
	}
}

func SmallStarsPublish() {
	for _, u := range SmallStars {
		for i := 0; i < NumSmallVideo; i++ {
			Publish(rand.Intn(NumUrls), int(u.Uid), NameSmallVideo)
		}

		var user User
		err := db.Model(&user).Where("id = ?", u.Uid).Update("work_count", NumSmallVideo).Error
		assert(err)
	}
	ch <- 1
	fmt.Println("Small stars Published over")
}

func MidStarsPublish() {
	for _, u := range MidStars {
		for i := 0; i < NumMidVideo; i++ {
			Publish(rand.Intn(NumUrls), int(u.Uid), NameMidVideo)
		}

		var user User
		err := db.Model(&user).Where("id = ?", u.Uid).Update("work_count", NumMidVideo).Error
		assert(err)
	}
	ch <- 1
	fmt.Println("Mid stars Published over")
}

func BigStarsPublish() {
	for _, u := range BigStars {
		for i := 0; i < NumBigVideo; i++ {
			Publish(rand.Intn(NumUrls), int(u.Uid), NameBigVideo)
		}

		var user User
		err := db.Model(&user).Where("id = ?", u.Uid).Update("work_count", NumBigVideo).Error
		assert(err)
	}
	ch <- 1
	fmt.Println("Big stars Published over")
}

func Publish(idx int, uid int, title string) {
	v := Video{
		PublishTime: TimeFormat,
		AuthorId:    int(uid),
		PlayUrl:     PlayUrls[idx],
		CoverUrl:    CoverUrls[idx],
		Title:       title,
	}
	if err := db.Create(&v).Error; err != nil {
		assert(err)
	}
}

func GetUrls() {
	file, err := os.OpenFile("video.txt", os.O_RDONLY, 0666)
	assert(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		str := string(line)
		strs := strings.Split(str, ",")
		PlayUrls[i], CoverUrls[i] = strs[0], strs[1]
	}

	ch <- 1
	fmt.Println("Got Urls")
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
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
	}

	ch <- 1
	fmt.Println("Got small stars info")
}
