package main

import (
	"context"
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

type storeVal struct {
	Value string
	Sub   *sub_once.SubOnce
}

var (
	configStore = make(map[string]*storeVal, 0)

	lock = &sync.RWMutex{}
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	cc := h.Group("/cc", middlewares.AccessLog())

	// 获取
	cc.GET("/read", func(c context.Context, ctx *app.RequestContext) {
		var key string

		if k := ctx.Query("key"); k == "" {
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  -1,
				"msg": "param is nil err",
			})
		} else {
			key = k
		}

		lock.RLock()
		val, ok := configStore[key]
		lock.RUnlock()
		if !ok {
			val = &storeVal{
				Value: "",
				Sub:   sub_once.New(),
			}

			lock.Lock()
			configStore[key] = val
			lock.Unlock()
		}

		timeout := time.After(5 * time.Second)
		ch := val.Sub.Sub()
		for {
			select {
			case <-timeout:
				ctx.JSON(consts.StatusOK, utils.H{
					"st":  0,
					"msg": "not update happen.",
				})
				return
			case <-ch:
				close(ch)
				ctx.JSON(consts.StatusOK, utils.H{
					"st":  0,
					"msg": "",
					"val": val.Value,
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

		// 修改
		lock.RLock()
		val, ok := configStore[key]
		lock.RUnlock()
		if ok {
			lock.Lock()
			val.Value = value
			lock.Unlock()
			val.Sub.Pub()
			ctx.JSON(consts.StatusOK, utils.H{
				"st":  0,
				"msg": "kv already modified.",
			})
			return
		} else {
			val = &storeVal{
				Value: value,
				Sub:   sub_once.New(),
			}
			lock.Lock()
			configStore[key] = val
			lock.Unlock()
			val.Sub.Pub()

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
