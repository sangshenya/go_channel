package ymtb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-channel/channel_handle/util"
	"io/ioutil"
	"net/http"
)

const(
	URL = "https://publisher-api.deepleaper.com/goods"
)

func Base(getReq *util.ReqMsg, failFunc util.ReqFailFunc, reqFunc util.ReqFunc, noFunc util.ReqNoFunc, timeoutFunc util.ReqTimeoutFunc, noimgFunc util.ReqNoimgFunc, nourlFunc util.ReqNourlFunc) util.ResMsg {
	uid := getReq.Imei
	uidtype := "imei"
	os := "android"
	if len(uid) == 0 {
		uid = getReq.Mac
		uidtype = "mac"
	}
	if getReq.Os == "2" {
		uid = getReq.Idfa
		uidtype = "idfa"
		os = "ios"
	}

	network := "unknown"
	switch getReq.Network {
	case "wifi":
		network = "wifi"
	case "2g":
		network = "2G"
	case "3g":
		network = "3G"
	case "4g":
		network = "4G"
	default:
		network = "unknown"
	}

	pid := getReq.ChannelReq.Adid
	channelid := getReq.ChannelReq.Appid

	if len(pid) == 0 || len(channelid) == 0 {
		failFunc(getReq)
		return util.ResMsg{}
	}

	// pid=**&channelid=**

	postdata := adreq{
		Version:    "1",
		Id:         util.GetRandom(),
		Pid:        pid,
		Channel_id: channelid,
		User:       _user{
			Uid:            uid,
			Uid_type:       uidtype,
			Uid_encryption: "NA",
		},
		Device:     _device{
			Ipv4:        getReq.Ip,
			Device_type: "phone",
			Device_make: getReq.Vendor,
			Device_os:   os,
			Network:     network,
		},
	}

	ma, error := json.Marshal(postdata)
	if error != nil {
		failFunc(getReq)
		return util.ResMsg{}
	}

	req, err := http.NewRequest("POST", URL, bytes.NewReader(ma))
	if err != nil {
		return util.ResMsg{}
	}
	//req.Header.Set("X-Forwarded-For", getReq.Ip)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")


	resp, err := util.Client.Do(req)
	reqFunc(getReq)

	if err != nil {
		fmt.Println("请求错误",err)
		noFunc(getReq)
		return util.ResMsg{}
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		noFunc(getReq)
		return util.ResMsg{}
	}

	if resp.StatusCode != 200 {
		noFunc(getReq)
		return util.ResMsg{}
	}

	resData := &adres{}
	json.Unmarshal(data, resData)

	if resData.Status != 0 {
		noFunc(getReq)
		return util.ResMsg{}
	}

	ad := resData.Creative

	imgurl := ad.Img
	if len(ad.Img) == 0 && len(ad.Imgs) != 0 {
		imgurl = ad.Imgs[0]
	}

	if len(imgurl) == 0 {
		noFunc(getReq)
		return util.ResMsg{}
	}

	postData := util.ResMsg{
		Title:                   ad.Title,
		Content:                 ad.Title,
		ImageUrl:                imgurl,
		Uri:                     ad.Clk_url,
		Scheme:                  ad.Deeplink_url,
		ImageList:               ad.Imgs,
	}

	return postData
}