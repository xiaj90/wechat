package miniRedPacket

import (
	"encoding/xml"
	"fmt"
	"github.com/xiaj90/wechat/pay/config"
	"github.com/xiaj90/wechat/util"
	"strconv"
	"time"
)

var miniRedPacketGateway = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendminiprogramhb"

// MiniRedPacket struct extends context
type MiniRedPacket struct {
	*config.Config
}

// NewMiniRedPacket return an instance of refund package
func NewMiniRedPacket(cfg *config.Config) *MiniRedPacket {
	miniRedPacket := MiniRedPacket{cfg}
	return &miniRedPacket
}

//Params 调用参数
type Params struct {
	MchBillno   string
	SendName    string
	ReOpenid    string
	TotalAmount string
	TotalNum    string
	Wishing     string
	ActName     string
	Remark      string
	NotifyWay   string
	SceneId     string
	RootCa      string //ca证书
}

//request 接口请求参数
type request struct {
	WxAppID     string `xml:"wxappid"`
	MchID       string `xml:"mch_id"`
	NonceStr    string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	MchBillno   string `xml:"mch_billno"`
	SendName    string `xml:"send_name"`
	ReOpenid    string `xml:"re_openid"`
	TotalAmount string `xml:"total_amount"`
	TotalNum    string `xml:"total_num"`
	Wishing     string `xml:"wishing"`
	ActName     string `xml:"act_name"`
	Remark      string `xml:"remark"`
	NotifyWay   string `xml:"notify_way"`
	SceneId     string `xml:"scene_id,omitempty"`
}

//Response 接口返回
type Response struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg"`
	WxAppID     string `xml:"wxappid,omitempty"`
	MchID       string `xml:"mch_id,omitempty"`
	ResultCode  string `xml:"result_code,omitempty"`
	ErrCode     string `xml:"err_code,omitempty"`
	ErrCodeDes  string `xml:"err_code_des,omitempty"`
	MchBillno   string `xml:"mch_billno,omitempty"`
	ReOpenid    string `xml:"re_openid,omitempty"`
	TotalAmount string `xml:"total_amount,omitempty"`
	SendListid  string `xml:"send_listid,omitempty"`
	Package     string `xml:"package,omitempty"`
}

//MiniRedPacket 发送红包
func (miniRedPacket *MiniRedPacket) RedPacket(p *Params) (rsp Response, err error) {
	nonceStr := util.RandomStr(32)
	param := make(map[string]string)
	param["wxappid"] = miniRedPacket.AppID
	param["mch_id"] = miniRedPacket.MchID
	param["nonce_str"] = nonceStr
	param["mch_billno"] = p.MchBillno
	param["send_name"] = p.SendName
	param["re_openid"] = p.ReOpenid
	param["total_amount"] = p.TotalAmount
	param["total_num"] = p.TotalNum
	param["wishing"] = p.Wishing
	param["act_name"] = p.ActName
	param["scene_id"] = p.SceneId
	param["remark"] = p.Remark
	param["notify_way"] = p.NotifyWay
	//param["sign_type"] = util.SignTypeMD5

	sign, err := util.ParamSign(param, miniRedPacket.Key)
	if err != nil {
		return
	}

	request := request{
		WxAppID:     miniRedPacket.AppID,
		MchID:       miniRedPacket.MchID,
		NonceStr:    nonceStr,
		Sign:        sign,
		MchBillno:   p.MchBillno,
		SendName:    p.SendName,
		ReOpenid:    p.ReOpenid,
		TotalAmount: p.TotalAmount,
		TotalNum:    p.TotalNum,
		Wishing:     p.Wishing,
		ActName:     p.ActName,
		Remark:      p.Remark,
		NotifyWay:   p.NotifyWay,
		SceneId:     p.SceneId,
	}
	rawRet, err := util.PostXMLWithTLS(miniRedPacketGateway, request, p.RootCa, miniRedPacket.MchID)
	if err != nil {
		return
	}
	err = xml.Unmarshal(rawRet, &rsp)
	if err != nil {
		return
	}
	if rsp.ReturnCode == "SUCCESS" {
		if rsp.ResultCode == "SUCCESS" {
			err = nil
			return
		}
		err = fmt.Errorf("refund error, errcode=%s,errmsg=%s", rsp.ErrCode, rsp.ErrCodeDes)
		return
	}
	err = fmt.Errorf("[msg : xmlUnmarshalError] [rawReturn : %s] [sign : %s]", string(rawRet), sign)
	return
}

func (miniRedPacket *MiniRedPacket) GetJsRedPacket(pkg string) (map[string]string, error) {
	var err error
	nonceStr := util.RandomStr(32)
	timeStamp := time.Now().Unix()
	param := make(map[string]string)
	param["appId"] = miniRedPacket.AppID
	param["nonceStr"] = nonceStr
	param["timeStamp"] = strconv.Itoa(int(timeStamp))
	param["package"] = pkg
	param["paySign"], err = util.ParamSign(param, miniRedPacket.Key)
	if err != nil {
		return nil, err
	}
	param["signType"] = util.SignTypeMD5
	return param, nil
}
