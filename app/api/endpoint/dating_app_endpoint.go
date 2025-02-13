package endpoint

import (
	"context"
	"dating-apps/app/api/middleware"
	"dating-apps/app/model/request"
	"dating-apps/app/service"

	"dating-apps/app/model/base"

	"github.com/go-kit/kit/endpoint"
)

type DatingAppEndpoint struct {
	SignUp          endpoint.Endpoint
	Login           endpoint.Endpoint
	Swipe           endpoint.Endpoint
	PurchasePremium endpoint.Endpoint
}

func MakeDatingAppEndpoint(s service.DatingAppService) DatingAppEndpoint {
	return DatingAppEndpoint{
		SignUp:          makeSignUpEndpoint(s),
		Login:           makeLoginEndpoint(s),
		Swipe:           makeSwipeEndpoint(s),
		PurchasePremium: makePurchasePremiumEndpoint(s),
	}
}

func makeSignUpEndpoint(s service.DatingAppService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (interface{}, error) {
		req := rqst.(request.SignUpRequest)
		result, msg := s.SignUp(ctx, req)
		return base.SetHttpResponse(ctx, msg, result), nil
	}
}

func makeLoginEndpoint(s service.DatingAppService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (interface{}, error) {
		req := rqst.(request.LoginRequest)
		result, msg := s.Login(ctx, req)
		return base.SetHttpResponse(ctx, msg, result), nil
	}
}

func makeSwipeEndpoint(s service.DatingAppService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (interface{}, error) {
		req := rqst.(request.SwipeRequest)
		// Extract the userID from the context.
		userIDUint, ok := ctx.Value(middleware.UserIDKey).(uint)
		if !ok {
			return nil, nil
		}

		req.UserID = userIDUint
		result, msg := s.Swipe(ctx, req)
		return base.SetHttpResponse(ctx, msg, result), nil
	}
}

func makePurchasePremiumEndpoint(s service.DatingAppService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (interface{}, error) {
		req := rqst.(request.PurchasePremiumRequest)
		// Extract the userID from the context.
		userIDUint, ok := ctx.Value(middleware.UserIDKey).(uint)
		if !ok {
			return nil, nil
		}

		req.UserID = userIDUint
		result, msg := s.PurchasePremium(ctx, req)
		return base.SetHttpResponse(ctx, msg, result), nil
	}
}
