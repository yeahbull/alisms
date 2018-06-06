# 阿里云短信服务SDK for Go

[![GoDoc](https://godoc.org/github.com/go-with/alisms?status.svg)](https://godoc.org/github.com/go-with/alisms)
[![Build Status](https://travis-ci.org/go-with/alisms.svg?branch=master)](https://travis-ci.org/go-with/alisms)

这是一个使用Go语言编写的阿里云短信服务SDK

## 举个栗子

以发送模板短信为栗

```Go
package main

import (
	"log"

	"github.com/go-with/alisms"
)

const (
	accessKeyID  = ""
	accessKeySecret  = ""
	signName = "" // 短信签名
	templateCode = "" // 短信模板CODE
	phoneNumbers = "" // 短信接收号码
)

func main() {
	// 实例化阿里云短信服务客户端
	c := alisms.NewClient(accessKeyID, accessKeySecret)
  
	// 发送模板短信
	_, err := c.SendSMS(phoneNumbers, signName, templateCode, alisms.TemplateParam{"code": "123456"})
	if err != nil {
		log.Fatal(err)
	}
}

```

## 许可协议

[The MIT License (MIT)](LICENSE)
