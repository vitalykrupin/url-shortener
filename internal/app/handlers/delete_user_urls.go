package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
)

type NewDeleteUserURLs struct {
	BaseHandler
}

func NewDeleteHandler(app *app.App) *NewDeleteUserURLs {
	return &NewDeleteUserURLs{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *NewDeleteUserURLs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		log.Println("Only DELETE requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userUUIDAny := req.Context().Value(middleware.UserIDKey)
	if userUUIDAny == nil || userUUIDAny == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID := userUUIDAny.(string)

	var urls []string
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Can not read body")
		return
	}
	if err := json.Unmarshal(body, &urls); err != nil {
		log.Println("Can not unmarshal body")
		return
	}

	out := handler.app.DeleteService.Add(urls)
	for s := range out {
		handler.app.Storage.DeleteUserURLs(req.Context(), userUUID, s)
	}
	// go func(ctx context.Context, in ...<-chan []string) {
	// 	handler.app.Storage.DeleteUserURLs(ctx, userUUID, urls)
	// }(req.Context(), out)

	// go func(ctx context.Context) {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		handler.app.Storage.DeleteUserURLs(ctx, userUUID, urls)
	// 		return
	// 	}
	// }(ctx)

	w.WriteHeader(http.StatusAccepted)
}
