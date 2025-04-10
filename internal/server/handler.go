package server

import (
	"net/http"
)

func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	for _, entry := range c.Cache {
		w.Write([]byte(entry.Metrics + "\n"))
	}
}
