package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating-apps/app"
	"dating-apps/app/api/endpoint"
	"dating-apps/app/model/base"
	"dating-apps/app/model/request"
	"dating-apps/app/service"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func DatingAppHttpHandler(s service.DatingAppService, app *app.Infra) http.Handler {
	route := mux.NewRouter()

	ep := endpoint.MakeDatingAppEndpoint(s)

	route.Methods(http.MethodPost).Path(app.UrlWithPrefix("user/swipe")).Handler(httptransport.NewServer(
		ep.Swipe,
		decodeSwipe,
		base.EncodeResponseHTTP,
	))

	return route
}

func decodeSwipe(ctx context.Context, r *http.Request) (interface{}, error) {
	var req request.SwipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}
