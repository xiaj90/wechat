package pay

import (
	"github.com/xiaj90/wechat/pay/config"
	"github.com/xiaj90/wechat/pay/miniRedPacket"
	"github.com/xiaj90/wechat/pay/notify"
	"github.com/xiaj90/wechat/pay/order"
	"github.com/xiaj90/wechat/pay/redPacket"
	"github.com/xiaj90/wechat/pay/refund"
	"github.com/xiaj90/wechat/pay/transfer"
)

//Pay 微信支付相关API
type Pay struct {
	cfg *config.Config
}

//NewPay 实例化微信支付相关API
func NewPay(cfg *config.Config) *Pay {
	return &Pay{cfg}
}

// GetOrder  下单
func (pay *Pay) GetOrder() *order.Order {
	return order.NewOrder(pay.cfg)
}

// GetNotify  通知
func (pay *Pay) GetNotify() *notify.Notify {
	return notify.NewNotify(pay.cfg)
}

// GetRefund  退款
func (pay *Pay) GetRefund() *refund.Refund {
	return refund.NewRefund(pay.cfg)
}

// GetTransfer  企业付款
func (pay *Pay) GetTransfer() *transfer.Transfer {
	return transfer.NewTransfer(pay.cfg)
}

// GetRedPacket  发送红包
func (pay *Pay) GetRedPacket() *redPacket.RedPacket {
	return redPacket.NewRedPacket(pay.cfg)
}

// GetMiniRedPacket  发送红包
func (pay *Pay) GetMiniRedPacket() *miniRedPacket.MiniRedPacket {
	return miniRedPacket.NewMiniRedPacket(pay.cfg)
}
