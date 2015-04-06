// +build !appengine

package bot

import (
	"log"
	"net/http"
)

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Client: new(http.Client),
		W:      w,
		R:      r,
	}
}

func (c *Context) logInfo(s string) {
	log.Println("INFO:", s)
}

func (c *Context) logDebug(s string) {
	log.Println("DEBUG:", s)
}

func (c *Context) logError(s string) {
	log.Println("ERROR:", s)
}
