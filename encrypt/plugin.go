package encrypt

import (
	"github.com/zzy-rabbit/bp/encrypt/internal"
	"github.com/zzy-rabbit/xtools/xcontext"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xplugin"
)

func init() {
	ctx := xcontext.Background()
	plugin := internal.New(ctx)
	err := xplugin.Register(ctx, plugin)
	if xerror.Error(err) {
		panic(err)
	}
}
