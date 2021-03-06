package dongqiudi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/sangshenya/go-channel/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const(
	URL_TEST = "http://test-union-ap.dongqiudi.com/ad-union/ads"
	URL = "http://union-ap.dongqiudi.com/ad-union/ads"

	URL_PRO = "http://113.209.195.3:5188/s2s/ds1/v1/sd5xo0t"

)

func Base(getReq *util.ReqMsg, reqFunc util.ReqFunc) (util.ResMsg, util.ChannelErrorProtocol) {

	if len(getReq.Carrier) == 0 {
		channelError := util.NewChannelRequestFailErrorWithText("请求参数运营商信息缺失")
		return util.ResMsg{}, channelError
	}

	network := 0
	switch getReq.Network {
	case "以太网":
		network = 1
	case "wifi":
		network = 2
	case "2g":
		network = 4
	case "3g":
		network = 6
	case "4g":
		network = 6
	}

	sd, err := strconv.Atoi(getReq.Sd)
	dpi := sd / 160
	den := float64(sd) / 160.0


	ostype := "android"
	if getReq.Os == "2" {
		ostype = "ios"
	}

	lat, err := strconv.ParseFloat(getReq.Lat, 64)
	if err != nil {
		lat = 0
	}
	lon, err := strconv.ParseFloat(getReq.Lng, 64)
	if err != nil {
		lon = 0
	}
	//
	swidth, err := strconv.Atoi(getReq.Screenwidth)
	if err != nil {
		swidth = 0
	}
	sheight, err := strconv.Atoi(getReq.Screenheight)
	if err != nil {
		sheight = 0
	}

	paramsMap := util.ParamsEncode(getReq.ChannelReq.Params, getReq.ChannelReq.Adtype)

	adid := getReq.ChannelReq.Adid
	appid := getReq.ChannelReq.Appid
	appname := getReq.ChannelReq.Appname
	pkg := getReq.ChannelReq.Pkg
	adtype := getReq.ChannelReq.Adtype
	w, _ := paramsMap["w"]
	h, _ := paramsMap["h"]

	if len(appid) == 0 || len(adid) == 0 || len(pkg) == 0 || len(adtype) == 0 {
		channelError := util.NewChannelRequestFailErrorWithText("请求必需参数部分参数为空")
		return util.ResMsg{}, channelError
	}

	postData := adreq{
		Id:util.Md5(util.GetRandom()),
		Version:"1.0.4",
		Imps:[]_imp{

		},
		App:_app{
			Id:       appid,
			Name:     appname,
			Bundle:   pkg,
			Ver:      getReq.Appversion,
			Paid:     0,
		},
		Device:_device{
			Ua:             getReq.Ua,
			Geo:            _geo{
				Lat: lat,
				Lon: lon,
			},
			Ip:             getReq.Ip,
			Ipv6:           getReq.Ip,
			DeviceType:     4,
			Make:           getReq.Vendor,
			Model:          getReq.Model,
			Os:             ostype,
			Osv:            getReq.Osversion,
			Rvs: 			getReq.Romversion,
			Sct:			getReq.Comptime,
			Anal: 			getReq.AndroidApiLevel,
			Hwv:            "",
			Sw:             swidth,
			Sh:             sheight,
			Ppi:            sd,
			Density: 		den,
			Dpi:            dpi,
			Ifa:            getReq.Idfa,
			Ifv:            getReq.Idfv,
			Did:            getReq.Imei,
			Dpid:           getReq.Androidid,
			Oaid:           getReq.Oaid,
			Mac:            getReq.Mac,
			Carrier:        getReq.Carrier,
			ConnectionType: network,
			Ibis:           getReq.Imsi_long,
			Orientation:    0,
			Language:       "zh",
		},
		At:1,
		Test:0,
		TMax:600,
		Ext:_ext{
			Rdt:      -1,
			Https:    -1,
			DeepLink: -1,
			Download: 1,
			Admt:     0,
			Vech:     0,
			Vecv:     0,
		},
		Language:"zh-CN",
	}

	aw, _ := strconv.Atoi(w)
	ah, _ := strconv.Atoi(h)

	// 1:banner 3:开屏 4:信息流
	imp := _imp{
		Id:       "1",
		Aw:       aw,
		Ah:       ah,
		TagId:    adid,
		BidFloor: 0,
		Banner:   nil,
		Native:   nil,
		Mts:      nil,
	}

	var mim = []string{"image/jpg","image/png","image/jpeg"}
	if adtype == "flow" {//flow
		imp.BidFloor = 100

		imp.Native = &_native{
			Assets:[]_assets{
				_assets{
					Id:1,
					Title:&_title{
						Len:20,
					},
				},
				_assets{
					Id:2,
					Data:&_data{
						Len:20,
					},
				},
				_assets{
					Id:3,
					Required:1,
					Img:&_img{
						W:aw,
						H:ah,
						Mimes:mim,
					},
				},
			},
		}
	} else {//banner startup
		if adtype == "banner" {
			imp.BidFloor = 150
		} else {
			imp.BidFloor = 700
		}
		imp.Banner = &_banner{
			W:aw,
			H:ah,
			Pos:0,
			Type:3,
			Mimes:mim,
		}
	}

	postData.Imps = append(postData.Imps, imp)

	ma, err := json.Marshal(&postData)
	if err != nil {
		channelError := util.NewChannelRequestFailErrorError(err)
		return util.ResMsg{}, channelError
	}


	req, err := http.NewRequest("POST", URL_PRO, bytes.NewReader(ma))
	if err != nil {
		channelError := util.NewChannelRequestFailErrorError(err)
		return util.ResMsg{}, channelError
	}
	req.Header.Set("User-Agent", getReq.Ua)
	req.Header.Set("X-Forwarded-For", getReq.Ip)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	resp, err := util.Client.Do(req)
	reqFunc(getReq)
	if err != nil {
		channelError := util.NewChannelRequestTimeoutError(err)
		return util.ResMsg{}, channelError
	}

	data := []byte{}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			channelError := util.NewChannelRequestNoError(err)
			return util.ResMsg{}, channelError
		}
		defer r.Close()

		data, err = ioutil.ReadAll(r)
		if err != nil {
			channelError := util.NewChannelRequestNoError(err)
			return util.ResMsg{}, channelError
		}
		resp.Body.Close()
	} else {
		data, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		channelError := util.NewChannelRequestNoError(err)
		return util.ResMsg{}, channelError
	}

	if resp.StatusCode != 200 {
		code := resp.StatusCode
		channelError := util.NewChannelRequestNoErrorWithText("状态码为:"+ strconv.Itoa(int(code)))
		return util.ResMsg{}, channelError
	}

	resData := &adres{}
	err = json.Unmarshal(data, resData)
	if err != nil {
		channelError := util.NewChannelRequestNoError(err)
		return util.ResMsg{}, channelError
	}

	if len(resData.Seats) == 0 {
		channelError := util.NewChannelRequestNoErrorWithText("Seats长度为0")
		return util.ResMsg{}, channelError
	}

	seats := resData.Seats[0]

	if len(seats.Bids) == 0 {
		channelError := util.NewChannelRequestNoErrorWithText("Bids长度为0")
		return util.ResMsg{}, channelError
	}

	ad := seats.Bids[0]

	title := ""
	content := ""

	imgUrl := ""
	clickUrl := ad.Target

	impArr := ad.Events.Els
	clkArr := ad.Events.Cls

	//fmt.Println(ad)

	//flow
	if adtype == "flow" {
		if len(ad.Native.Assets) == 0 {
			channelError := util.NewChannelRequestNoErrorWithText("ad.Native.Assets长度为0")
			return util.ResMsg{}, channelError
		}
		for _,item := range ad.Native.Assets {
			itemImg := item.Img
			itemTitle := item.Title
			itemDec := item.Data
			if len(itemImg.Url) != 0 {
				imgUrl = itemImg.Url
			}
			if len(itemTitle.Text) != 0 {
				title = itemTitle.Text
			}
			if len(itemDec.Value) != 0 {
				content = itemDec.Value
			}

		}
	} else {
		//fmt.Println(ad.Banner.Iurl)
		imgUrl = ad.Banner.Url

	}

	if len(imgUrl) == 0 {
		channelError := util.NewChannelNoImageErrorWithText("图片链接长度为0")
		return util.ResMsg{}, channelError
	}

	pos := util.CreateAbScreenPosWHWithOs(getReq.Screenwidth, getReq.Screenheight, adtype, getReq.Os)

	resultData := util.ResMsg{
		Id:       util.Md5(string(data) + time.Now().String()),
		Weight:   0,
		State:    0,
		Title:    title,
		Content:  content,
		ImageUrl: imgUrl,
		Uri:      replace(clickUrl, pos),
	}

	if len(ad.DeepLink) != 0 {
		resultData.Scheme = ad.DeepLink
		resultData.Schemereport = ad.Events.Dcls
	}

	if ad.ActionType == "2" && ad.Demand == "GDT" {
		resultData.Json = true
	}

	displayArray := []string{}
	clickArray := []string{}
	for _, item := range impArr {
		displayArray = append(displayArray, replace(item, pos))
	}

	for _, item := range clkArr {
		clickArray = append(clickArray, replace(item, pos))
	}

	resultData.Displayreport = displayArray
	resultData.Clickreport = clickArray

	for _, item := range ad.Events.Sdls {
		resultData.StartDownload = append(resultData.StartDownload, replace(item, pos))
	}

	for _, item := range ad.Events.Edls {
		resultData.Downloaded = append(resultData.Downloaded, replace(item, pos))
	}

	for _, item := range ad.Events.Sils {
		resultData.Installed = append(resultData.Installed, replace(item, pos))
	}

	for _, item := range ad.Events.Eils {
		resultData.Installed = append(resultData.Installed, replace(item, pos))
	}

	for _, item := range ad.Events.Ials {
		resultData.Open = append(resultData.Open, replace(item, pos))
	}

	if resultData.ResponseDataIsEmpty(getReq.Adtype) {
		channelError := util.NewChannelRequestNoErrorWithText("数据不完整")
		return util.ResMsg{}, channelError
	}

	return resultData, nil

}

func replace(urlStr string, pos [10]string) string {

	urlStr = strings.Replace(urlStr, "__TS__", util.TS, -1)
	urlStr = strings.Replace(urlStr, "__AZMTS__", util.TS, -1)
	urlStr = strings.Replace(urlStr, "__STS__", util.TS, -1)
	urlStr = strings.Replace(urlStr, "__ETS__", util.TS, -1)
	urlStr = strings.Replace(urlStr, "__AZCTS__", util.TS, -1)

	urlStr = strings.Replace(urlStr, "__AZCX__", util.DX, -1)
	urlStr = strings.Replace(urlStr, "__AZCY__", util.DY, -1)
	urlStr = strings.Replace(urlStr, "__AZMX__", util.UX, -1)
	urlStr = strings.Replace(urlStr, "__AZMY__", util.UY, -1)

	urlStr = strings.Replace(urlStr, "__DSMX__", util.RDX, -1)
	urlStr = strings.Replace(urlStr, "__DSMY__", util.RDY, -1)
	urlStr = strings.Replace(urlStr, "__DSCX__", util.RUX, -1)
	urlStr = strings.Replace(urlStr, "__DSCY__", util.RUY, -1)


	urlStr = strings.Replace(urlStr, "__AMVW__", pos[8], -1)
	urlStr = strings.Replace(urlStr, "__AMVH__", pos[9], -1)

	urlStr = strings.Replace(urlStr, "__AMSW__", pos[8], -1)
	urlStr = strings.Replace(urlStr, "__AMSH__", pos[9], -1)

	return urlStr
}