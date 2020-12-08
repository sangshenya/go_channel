package test

import (
	"github.com/sangshenya/go-channel/util"
	"time"
)

type adInfo struct {
	Title string
	Content string
	ImageUrl string
}
// "",""
func Base(getReq *util.ReqMsg, failFunc util.ReqFailFunc, reqFunc util.ReqFunc, noFunc util.ReqNoFunc, timeoutFunc util.ReqTimeoutFunc, noimgFunc util.ReqNoimgFunc, nourlFunc util.ReqNourlFunc) util.ResMsg {
	reqFunc(getReq)
	adinfo := adInfo{
		Title:    "大牌好货,每满300减40立即前往",
		Content:  "每满300减40立即前往，普通测试广告",
		ImageUrl: "https://img.admobile.top/admobile-adRequest/dhh_1.png",
	}

	resultData := util.ResMsg{
		Id:       util.Md5(util.GetRandom() + time.Now().String()),
		Weight:   0,
		State:    0,
		Title:    adinfo.Title,
		Content:  adinfo.Content,
		ImageUrl: adinfo.ImageUrl,
		Uri:      "https://pages.tmall.com/wow/a/act/tmall/dailygroup/1773/wupr?wh_pid=daily-218058&activity_id=100000000145",
	}

	switch getReq.ChannelReq.Adtype {
		case "flow":
			resultData.ImageUrl = "https://admobile.oss-cn-hangzhou.aliyuncs.com/admobile-adRequest/tbdhh_fz.jpg"
		case "banner":
			resultData.ImageUrl = "https://admobile.oss-cn-hangzhou.aliyuncs.com/admobile-adRequest/71911590995175_.pic_hd.jpg"
	}

	if len(resultData.ImageUrl) == 0 {
		noimgFunc(getReq)
		return util.ResMsg{}
	}

	if len(resultData.Uri) == 0 {
		nourlFunc(getReq)
		return util.ResMsg{}
	}

	return resultData
}