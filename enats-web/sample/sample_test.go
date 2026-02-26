package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"log"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// 假设这是生成的 Protobuf 结构
type Person struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age   int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	//CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,proto3" json:"created_at,omitempty"`
}

func TestName(t *testing.T) {
	//dynamicWithMap()

	person1 := &Person{
		Name:  "李四",
		Age:   30,
		Email: "lisi@example.com",
	}

	m := gjson.New(person1).Map()
	person, err := structpb.NewStruct(m)
	if err != nil {
		log.Fatal(err)
	}
	// 1. 二进制格式（最紧凑）
	binaryData, _ := proto.Marshal(person)
	fmt.Printf("二进制大小: %d 字节\n", len(binaryData))

	// 2. Text 格式（可读）
	textData := prototext.Format(person)
	fmt.Printf("Text 格式: %s\n", textData)
	fmt.Printf("Text 大小: %d 字节\n", len(textData))

	// 3. JSON 格式
	jsonData, _ := protojson.Marshal(person)
	fmt.Printf("JSON 格式: %s\n", jsonData)
	fmt.Printf("JSON 大小: %d 字节\n", len(jsonData))

	// 4. 压缩后的二进制
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	gz.Write(binaryData)
	gz.Close()
	fmt.Printf("压缩后大小: %d 字节\n", compressed.Len())
	fmt.Printf("压缩率: %.2f%%\n",
		100-float64(compressed.Len())/float64(len(binaryData))*100)
}

func dynamicWithMap() {
	jsonStr := `{"address":{"city":"北京", "street":"长安街"}, "age":25, "name":"张三", "skills":["Go", "Python"]}`

	// 1. 先解析为 map
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		log.Fatal(err)
	}

	// 2. 动态处理数据
	fmt.Printf("姓名: %v\n", data["name"])
	fmt.Printf("年龄: %v\n", data["age"])

	// 3. 转换为 Struct (Protobuf 的动态结构)
	structValue, err := structpb.NewStruct(data)
	if err != nil {
		log.Fatal(err)
	}

	s := structpb.Struct{}

	marshal, err := proto.Marshal(structValue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Struct 长度:%d \n", len(marshal))
	fmt.Printf("Struct 转换为二进制: %x\n", marshal)

	err = proto.Unmarshal(marshal, &s)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Struct 转 JSON
	jsonBytes, _ := protojson.Marshal(&s)
	fmt.Printf("Struct 转 JSON: %s\n", jsonBytes)
	fmt.Printf(" JSON 长度:  %d\n", len(jsonBytes))
}
