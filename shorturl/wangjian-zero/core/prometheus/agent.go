package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"shorturl/wangjian-zero/core/logx"
	"shorturl/wangjian-zero/core/threading"
)

var once sync.Once

// StartAgent starts a prometheus agent.
func StartAgent(c Config) {
	once.Do(func() {
		if len(c.Host) == 0 {
			return
		}

		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			logx.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logx.Error(err)
			}
		})
	})
}
