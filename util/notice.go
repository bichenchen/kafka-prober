package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kafka-prober/logger"
	"net"
	"net/http"
	"time"
)

type AlertPck struct {
	Hi_channel  string `json:"channel"`
	Receiver    string `json:"receiver"`
	Description string `json:"description"`
}

type HiPusher struct {
	AppId     string      `json:"appId"`
	Token     string      `json:"token"`
	AlertInfo [1]AlertPck `json:"alertList"`
}

const (
	HOST            = "http://1.1.1.1"
	URL             = "/alert/push"
	APP_ID          = "111"
	TOKEN           = "111"
	CHANNEL_TYPE_HI = "group"
	RECEIVER        = "111"
)

func sendPostRequest(client *http.Client, host string, uri string, sendObj HiPusher) ([]byte, error) {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(sendObj)
	// fmt.Println(requestBody)
	req, err := http.NewRequest("POST", host+uri, requestBody)
	// fmt.Println(req)
	if err != nil {
		logger.Logger.Warnf("new request failed - %s", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		logger.Logger.Warnf("Failed to send msg to Hermes %s", err)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Warnf("Failed to parse result from Hermes, err=%v", err)
		return nil, err
	}

	return result, nil
}

func SendAlertInfoToHi(clusterName string, brokers string, mode string,
	topic string, groupID string, message string) {
	alt := AlertPck{
		Hi_channel: CHANNEL_TYPE_HI,
		Receiver:   RECEIVER,
	}
	sendObj := HiPusher{
		AppId: APP_ID,
		Token: TOKEN,
	}
	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	var hiMsg string
	if mode == "produce" {
		hiMsg = fmt.Sprintf("【生产者模式】\n【Time】: %s\n【Cluster】: %s\n【Brokers】: %s\n【Topic】: %s\n【Message】: %s",
			nowFormat, clusterName, brokers, topic, message)
	} else {
		hiMsg = fmt.Sprintf("【消费者模式】\n【Time】: %s\n【Cluster】: %s\n【Brokers】: %s\n【Topic】: %s\n【GroupID】: %s\n【Message】: %s",
			nowFormat, clusterName, brokers, topic, groupID, message)
	}
	sendObj.AlertInfo[0] = alt
	sendObj.AlertInfo[0].Description = hiMsg

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				// set connect timeoue
				conn, err := net.DialTimeout(netw, addr, time.Millisecond*500)

				if err != nil {
					return nil, err
				}

				// set read/write timeout
				conn.SetDeadline(time.Now().Add(time.Millisecond * 500))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	_, err := sendPostRequest(client, HOST, URL, sendObj)
	if err != nil {
		logger.Logger.Warnf("send message error:%s", err)
	}
	// fmt.Printf("send message result:%s", resp)
}
