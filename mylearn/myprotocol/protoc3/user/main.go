package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"mylearn/myprotocol/protoc3/user/user"
)

func main() {
	article := &user.Article{
		Aid:   1,
		Title: "protobuf for golang",
		Views: 100,
	}
	//序列化成二进制数据
	bytes, _ := proto.Marshal(article)
	fmt.Printf("bytes: %v\n", bytes)
	//反序列化
	otherArticle := &user.Article{}
	proto.Unmarshal(bytes, otherArticle)
	fmt.Printf("otherArticle.GetAid():%v\n", otherArticle.GetAid())
	fmt.Printf("otherArticle.GetTitle():%v\n", otherArticle.GetTitle())
	fmt.Printf("otherArticle.GetViews():%v\n", otherArticle.GetViews())
}
