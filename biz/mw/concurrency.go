package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"sync"
)

func MaxWorker(max int, skipper ...SkipperFunc) app.HandlerFunc {
	hlog.Infof("setting max workers to %d", max)
	sem := make(chan struct{}, max)
	return func(ctx context.Context, c *app.RequestContext) {
		sem <- struct{}{}
		defer func() {
			<-sem
		}()
		c.Next(ctx)
	}
}

type MaxRequestIface struct {
	current int
	lock    *sync.RWMutex
}

func MaxRequest(max int, skipper ...SkipperFunc) app.HandlerFunc {
	hlog.Infof("setting max requests to %d", max)
	m := &MaxRequestIface{
		current: 0,
		lock:    &sync.RWMutex{},
	}

	return func(ctx context.Context, c *app.RequestContext) {
		m.lock.RLock()
		if m.current >= max {
			m.lock.RUnlock()
			c.JSON(http.StatusServiceUnavailable, map[string]any{"code": -503, "message": "Too many requests"})
			c.Abort()
			return
		}
		m.lock.RUnlock()
		m.lock.Lock()
		m.current++
		m.lock.Unlock()
		c.Next(ctx)
		m.lock.Lock()
		m.current--
		m.lock.Unlock()
	}
}
