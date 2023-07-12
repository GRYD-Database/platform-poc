package server

import "net/http"

func (c *Container) grydAccessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.grydSemaphore.TryAcquire(1) {
			c.logger.Debug("gryd access: simultaneous on-chain operations not supported")
			c.logger.Error(nil, "staking access: simultaneous on-chain operations not supported")
			WriteJson(w, "simultaneous on-chain operations not supported", http.StatusTooManyRequests)
			return
		}
		defer c.grydSemaphore.Release(1)
	})
}
