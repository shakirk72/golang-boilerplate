package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/helpers"
)

// Drainer tracks in-flight requests and blocks new ones once Drain is called.
// The mutex guards the shutdown flag and wg.Add together to prevent a race
// where a handler goroutine passes the shutdown check but hasn't called
// wg.Add yet when Drain sets the flag and calls wg.Wait.
type Drainer struct {
	mu       sync.Mutex
	wg       sync.WaitGroup
	shutdown bool
}

// Handler is a Gin middleware. Once Drain has been called it immediately
// returns 503 to any new request; otherwise it tracks the request in the
// WaitGroup for the duration of the handler chain.
func (d *Drainer) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		d.mu.Lock()
		if d.shutdown {
			d.mu.Unlock()
			c.JSON(http.StatusServiceUnavailable, helpers.Response{
				Success: false,
				Message: "server is shutting down",
			})
			c.Abort()
			return
		}
		d.wg.Add(1)
		d.mu.Unlock()

		defer d.wg.Done()
		c.Next()
	}
}

// Drain marks the server as shutting down and blocks until every in-flight
// request has completed. Call this before http.Server.Shutdown.
func (d *Drainer) Drain() {
	d.mu.Lock()
	d.shutdown = true
	d.mu.Unlock()
	d.wg.Wait()
}
