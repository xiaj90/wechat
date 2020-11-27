package notify

import (
	"github.com/xiaj90/wechat/pay/config"
)

//Notify 回调
type Notify struct {
	*config.Config
}

//NewNotify new
func NewNotify(cfg *config.Config) *Notify {
	return &Notify{cfg}
}
