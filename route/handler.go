package route

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//Handlers custom handlers
type Handlers []http.Handler

// Implement the ServerHTTP method on our new type
func (handles Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handler := range handles {
		if strings.HasPrefix(r.URL.Path, "/api/v") {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
		}
		handler.ServeHTTP(w, r)
	}
}

// InitHandlers for custom handlers
func InitHandlers() *Handlers {

	router := httprouter.New()
	router.GET("/api/v2/repository/tag/:id", GetTag)
	router.GET("/api/v2/tags", GetTags)
	router.GET("/api/v2/repositories", GetRepositories)

	mHandlers := make(Handlers, 0)
	mHandlers = append(mHandlers, router)
	log.Fatal(http.ListenAndServe(":8080", mHandlers))

	return &mHandlers
}
