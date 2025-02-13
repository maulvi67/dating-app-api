package handler

import (
	"context"
	"dating-apps/app"
	"dating-apps/app/api/endpoint"
	"dating-apps/app/model/base"
	"dating-apps/app/model/request"
	"dating-apps/app/service"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func AuthHttpHandler(s service.DatingAppService, app *app.Infra) http.Handler {
	route := mux.NewRouter()
	ep := endpoint.MakeDatingAppEndpoint(s)
	route.Methods(http.MethodPost).Path(app.UrlWithPrefix("auth/signup")).Handler(httptransport.NewServer(
		ep.SignUp,
		decodeSignUp,
		base.EncodeResponseHTTP,
	))

	route.Methods(http.MethodPost).Path(app.UrlWithPrefix("auth/login")).Handler(httptransport.NewServer(
		ep.Login,
		decodeLogin,
		base.EncodeResponseHTTP,
	))
	return route
}

func decodeSignUp(ctx context.Context, r *http.Request) (interface{}, error) {
	var req request.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLogin(ctx context.Context, r *http.Request) (interface{}, error) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}
