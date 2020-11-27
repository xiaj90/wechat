package context

import (
	"github.com/xiaj90/wechat/credential"
	"github.com/xiaj90/wechat/officialaccount/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
