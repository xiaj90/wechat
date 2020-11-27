package transfer

import (
	"encoding/xml"
	"fmt"
	"github.com/xiaj90/wechat/pay/config"
	"github.com/xiaj90/wechat/util"
)

var transferGateway = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"

// Transfer struct extends context
type Transfer struct {
	*config.Config
}

// NewTransfer return an instance of refund package
func NewTransfer(cfg *config.Config) *Transfer {
	transfer := Transfer{cfg}
	return &transfer
}

//Params 调用参数
type Params struct {
	DeviceInfo     string
	PartnerTradeNo string
	Openid         string
	CheckName      string
	ReUserName     string
	Amount         string
	Desc           string
	SpillCreateIp  string
	RootCa         string //ca证书
}

//request 接口请求参数
type request struct {
	MchAppID       string `xml:"mch_appid"`
	MchID          string `xml:"mchid"`
	NonceStr       string `xml:"nonce_str"`
	DeviceInfo     string `xml:"device_info"`
	Sign           string `xml:"sign"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	Openid         string `xml:"openid"`
	CheckName      string `xml:"check_name"`
	ReUserName     string `xml:"re_user_name"`
	Amount         string `xml:"amount"`
	Desc           string `xml:"desc"`
	SpillCreateIp  string `xml:"spbill_create_ip,omitempty"`
}

//Response 接口返回
type Response struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	MchAppID       string `xml:"mch_appid,omitempty"`
	MchID          string `xml:"mchid,omitempty"`
	DeviceInfo     string `xml:"device_info,omitempty"`
	NonceStr       string `xml:"nonce_str,omitempty"`
	ResultCode     string `xml:"result_code,omitempty"`
	ErrCode        string `xml:"err_code,omitempty"`
	ErrCodeDes     string `xml:"err_code_des,omitempty"`
	PartnerTradeNo string `xml:"partner_trade_no,omitempty"`
	PaymentNo      string `xml:"payment_no,omitempty"`
	PaymentTime    string `xml:"payment_time,omitempty"`
}

//Transfer 付款到零钱
func (transfer *Transfer) Transfer(p *Params) (rsp Response, err error) {
	nonceStr := util.RandomStr(32)
	param := make(map[string]string)
	param["mch_appid"] = transfer.AppID
	param["mchid"] = transfer.MchID
	param["nonce_str"] = nonceStr
	param["device_info"] = p.DeviceInfo
	param["partner_trade_no"] = p.PartnerTradeNo
	param["openid"] = p.Openid
	param["check_name"] = p.CheckName
	param["re_user_name"] = p.ReUserName
	param["amount"] = p.Amount
	param["desc"] = p.Desc
	param["spbill_create_ip"] = p.SpillCreateIp
	//param["sign_type"] = util.SignTypeMD5

	sign, err := util.ParamSign(param, transfer.Key)
	if err != nil {
		return
	}

	request := request{
		MchAppID:       transfer.AppID,
		MchID:          transfer.MchID,
		NonceStr:       nonceStr,
		Sign:           sign,
		DeviceInfo:     p.DeviceInfo,
		PartnerTradeNo: p.PartnerTradeNo,
		Openid:         p.Openid,
		CheckName:      p.CheckName,
		ReUserName:     p.ReUserName,
		Amount:         p.Amount,
		Desc:           p.Desc,
		SpillCreateIp:  p.SpillCreateIp,
	}
	rawRet, err := util.PostXMLWithTLS(transferGateway, request, p.RootCa, transfer.MchID)
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
