package main

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"mylearn/myprotocol/protoc2/blog"
)

func main() {
	article := &blog.Article{
		Aid:   1,
		Title: "protobuf for golang",
		Views: 100,
	}
	//Message to json
	jsonString := protojson.Format(article.ProtoReflect().Interface())
	fmt.Printf("jsonString: %v\n", jsonString)

	//json to Message
	m := article.ProtoReflect().Interface()
	protojson.Unmarshal([]byte(jsonString), m)

	fmt.Printf("m:%v/n", m)
}
