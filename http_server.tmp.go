package main

/*

package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/daidai21/long_polling_demo/pkg/middlewares"
	"github.com/daidai21/long_polling_demo/pkg/sub_once"
	"sync"
	"time"
)

// TODO: not test

var configStore = sync.Map{} // Key: string; Value: &storeValue

type (
	storeKey   string
	storeValue struct {
		Value *string
		Sub   *sub_once.SubOnce
	}
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	cc := h.Group("/cc", middlewares.AccessLog())

	// 获取
	cc.GET("/read", func(c context.Context, ctx *app.RequestContext) {
		var key storeKey

		if k := ctx.Query("key"); k == "" {
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  -1,
				"msg": "param is nil err",
			})
		} else {
			key = storeKey(k)
		}

		hlog.Debug("DEBUG read", key)

		val, ok := configStore.Load(key)
		if !ok {
			var value string
			val = storeValue{
				Value: &value,
				Sub:   sub_once.New(),
			}
			configStore.Store(key, val)
		}
		storeVal := val.(storeValue)

		timeout := time.After(5 * time.Second)
		ch := storeVal.Sub.Sub()
		for {
			select {
			case <-timeout:
				ctx.JSON(consts.StatusOK, utils.H{
					"st":  0,
					"msg": "not update happen.",
				})
				return
			case <-ch:

				fmt.Println("for case ch")

				close(ch)
				ctx.JSON(consts.StatusOK, utils.H{
					"st":  0,
					"msg": "key: " + *storeVal.Value,
				})
				return
			}
		}

	})

	// 创建或更新
	cc.GET("/write", func(c context.Context, ctx *app.RequestContext) {
		var key, value string

		if k := ctx.Query("key"); k == "" {
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  -1,
				"msg": "param key is nil err",
			})
			return
		} else {
			key = k
		}
		if v := ctx.Query("value"); v == "" {
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  -1,
				"msg": "param value is nil err",
			})
			return
		} else {
			value = v
		}

		hlog.Debug("DEBUG write", key, value)

		// 修改
		val, ok := configStore.Load(key)
		if ok {
			storeVal := val.(storeValue)
			storeVal.Value = &value
			storeVal.Sub.Pub()
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  0,
				"msg": "kv already modified.",
			})
			return
		} else {
			var value string
			val = storeValue{
				Value: &value,
				Sub:   sub_once.New(),
			}
			configStore.Store(key, val)
			storeVal := val.(storeValue)
			storeVal.Sub.Pub()

			ctx.JSON(consts.StatusOK, utils.H{
				"st":  0,
				"msg": "kv already created.",
			})
			return
		}

	})

	hlog.SetLevel(hlog.LevelDebug)
	hlog.Info(h.Routes()) // routeInfo

	h.Spin()
}


*/
