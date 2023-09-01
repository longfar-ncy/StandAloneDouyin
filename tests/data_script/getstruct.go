package main

import (
	"fmt"

	"github.com/gohouse/converter"
)

func main() {
	err := converter.NewTable2Struct().
		SavePath("./struct.go").
		Dsn("root:root@tcp(127.0.0.1:3306)/douyin").
		TagKey("gorm").
		EnableJsonTag(true).
		Table("user").
		Run()
	fmt.Println(err)
}
