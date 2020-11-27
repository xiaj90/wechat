package redPacket

import (
	"encoding/xml"
	"fmt"

	"github.com/xiaj90/wechat/pay/config"
	"github.com/xiaj90/wechat/util"
)

var redPacketGateway = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"

// Refund struct extends context
type RedPacket struct {
	*config.Config
}

// NewRedPacket return an instance of refund package
func NewRedPacket(cfg *config.Config) *RedPacket {
	redPacket := RedPacket{cfg}
	return &redPacket
}

//Params 调用参数
type Params struct {
	MchBillno   string
	SendName    string
	ReOpenid    string
	TotalAmount string
	TotalNum    string
	Wishing     string
	ClientIp    string
	ActName     string
	Remark      string
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
	ClientIp    string `xml:"client_ip"`
	ActName     string `xml:"act_name"`
	Remark      string `xml:"remark"`
	SceneId     string `xml:"scene_id,omitempty"`
	RiskInfo    string `xml:"risk_info,omitempty"`
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
}

//RedPacket 发送红包
func (redPacket *RedPacket) RedPacket(p *Params) (rsp Response, err error) {
	nonceStr := util.RandomStr(32)
	param := make(map[string]string)
	param["wxappid"] = redPacket.AppID
	param["mch_id"] = redPacket.MchID
	param["nonce_str"] = nonceStr
	param["mch_billno"] = p.MchBillno
	param["send_name"] = p.SendName
	param["re_openid"] = p.ReOpenid
	param["total_amount"] = p.TotalAmount
	param["total_num"] = p.TotalNum
	param["wishing"] = p.Wishing
	param["client_ip"] = p.ClientIp
	param["act_name"] = p.ActName
	param["scene_id"] = p.SceneId
	param["remark"] = p.Remark
	//param["sign_type"] = util.SignTypeMD5

	sign, err := util.ParamSign(param, redPacket.Key)
	if err != nil {
		return
	}

	request := request{
		WxAppID:     redPacket.AppID,
		MchID:       redPacket.MchID,
		NonceStr:    nonceStr,
		Sign:        sign,
		MchBillno:   p.MchBillno,
		SendName:    p.SendName,
		ReOpenid:    p.ReOpenid,
		TotalAmount: p.TotalAmount,
		TotalNum:    p.TotalNum,
		Wishing:     p.Wishing,
		ClientIp:    p.ClientIp,
		ActName:     p.ActName,
		Remark:      p.Remark,
		SceneId:     p.SceneId,
	}
	rawRet, err := util.PostXMLWithTLS(redPacketGateway, request, p.RootCa, redPacket.MchID)
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
