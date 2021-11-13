package server

import (
	"bytes"
	"crypto/des"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func convertMonthToInt(month string) int {
	var m map[string]int
	m = make(map[string]int)
	m["January"] = 1
	m["February"] = 2
	m["March"] = 3
	m["April"] = 4
	m["May"] = 5
	m["June"] = 6
	m["July"] = 7
	m["August"] = 8
	m["September"] = 9
	m["October"] = 10
	m["November"] = 11
	m["December"] = 12
	return m[month]
}


type NotifReqTelegram struct {
	ChatId              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}



func SendNotificationByTelegram(message string, title string) {
	log.Println("Send by Telegram started")
	url := "https://api.telegram.org/bot1908920066:AAH83I6JFKGsWfE1f20f0y_S-6NDHKEjWW4/sendMessage"
	jsonBytes, err := json.Marshal(&NotifReqTelegram{
		ChatId:              "-1001435126738",
		Text:                title + "\n" + message,
		DisableNotification: false,
	})
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
	}
	//defer resp.Body.Close()
	//
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	log.Println("Send by Telegram ended")
}


func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func DesEncrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	// src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func DesDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	// out = PKCS5UnPadding(out)
	return out, nil
}
