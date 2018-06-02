package alisms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/satori/go.uuid"
)

const (
	urlPref = "http://dysmsapi.aliyuncs.com/?"
)

var (
	specialURLReplacer = strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
)

type Client struct {
	AccessKeyID     string
	AccessKeySecret string
}

func New(accessKeyID, accessKeySecret string) *Client {
	return &Client{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	}
}

func (c *Client) SendSMS(phoneNumbers, signName, templateCode string, templateParam interface{}) (*simplejson.Json, error) {
	paras := make(map[string]string)

	// 系统参数
	paras["Action"] = "SendSms"
	paras["Version"] = "2017-05-25"
	paras["RegionId"] = "cn-hangzhou"
	paras["PhoneNumbers"] = phoneNumbers
	paras["SignName"] = signName
	paras["TemplateCode"] = templateCode
	if templateParam != nil {
		json, err := json.Marshal(templateParam)
		if err != nil {
			return nil, err
		}
		paras["TemplateParam"] = string(json)
	}

	// 业务参数
	loc, err := time.LoadLocation("GMT0")
	if err != nil {
		return nil, err
	}
	nonce, err := uuid.NewV1()
	if err != nil {
		return nil, err
	}
	paras["AccessKeyId"] = c.AccessKeyID
	paras["Timestamp"] = time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	paras["Format"] = "JSON"
	paras["SignatureMethod"] = "HMAC-SHA1"
	paras["SignatureVersion"] = "1.0"
	paras["SignatureNonce"] = nonce.String()

	// 字典排序参数键
	keys := make([]string, 0, len(paras))
	for k, _ := range paras {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构造查询字符串
	query := new(bytes.Buffer)
	for _, k := range keys {
		query.WriteString(specialURLEncode(k))
		query.WriteString("=")
		query.WriteString(specialURLEncode(paras[k]))
		query.WriteString("&")
	}
	signStr := fmt.Sprint("GET&", specialURLEncode("/"), "&", specialURLEncode(strings.TrimRight(query.String(), "&"))) // 构造签名字符串
	signature := c.sign(signStr)
	query.WriteString("Signature=")
	query.WriteString(specialURLEncode(signature))

	// 构造请求URL
	urlStr := fmt.Sprint(urlPref, query.String())

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (c *Client) sign(signStr string) string {
	mac := hmac.New(sha1.New, []byte(c.AccessKeySecret+"&"))
	mac.Write([]byte(signStr))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func specialURLEncode(s string) string {
	return specialURLReplacer.Replace(url.QueryEscape(s))
}
