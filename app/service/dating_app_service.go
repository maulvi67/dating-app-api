package service

import (
	"context"
	"dating-apps/app"
	"dating-apps/app/model/entity"
	"dating-apps/app/model/request"
	"dating-apps/app/model/response"
	"dating-apps/app/repository"
	"dating-apps/helper/message"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type DatingAppService interface {
	SignUp(ctx context.Context, input request.SignUpRequest) (*response.SignUpResponse, message.Message)
	Login(ctx context.Context, input request.LoginRequest) (*response.LoginResponse, message.Message)
	Swipe(ctx context.Context, input request.SwipeRequest) (*response.SwipeResponse, message.Message)
}

type datingAppService struct {
	app        *app.Infra
	repository repository.DatingAppRepository
}

func NewDatingAppService(
	app *app.Infra,
	repo repository.DatingAppRepository,
) DatingAppService {
	return &datingAppService{
		app:        app,
		repository: repo,
	}
}

func (s *datingAppService) SignUp(ctx context.Context, input request.SignUpRequest) (*response.SignUpResponse, message.Message) {
	// Check if the user already exists.
	count, err := s.repository.CountUser(input.Email)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}
	if count > 0 {
		return nil, message.ErrDataExists
	}

	// Hash the password.
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}

	user := entity.User{
		Email:        input.Email,
		Name:         input.Name,
		Gender:       input.Gender,
		PasswordHash: string(hashed),
		IsPremium:    false,
	}

	// Create the user.
	_, err = s.repository.CreateUser(user)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}

	// You can return additional user data if needed.
	return &response.SignUpResponse{
		Message: "User created successfully",
	}, message.SuccessMsg
}

func (s *datingAppService) Login(ctx context.Context, input request.LoginRequest) (*response.LoginResponse, message.Message) {
	// Retrieve the user by email.
	user, err := s.repository.GetUserByEmail(input.Email)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.ErrInvalidCred
	}

	// Compare provided password with stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.ErrInvalidCred
	}

	// Create JWT token.
	expirationTime := time.Now().Add(time.Duration(s.app.Config.SecurityConfig.JwtConfig.JwtExpireHours) * time.Hour)
	claims := &entity.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.app.Config.SecurityConfig.JwtConfig.JwtSecret))
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}

	return &response.LoginResponse{
		Token: tokenString,
	}, message.SuccessMsg
}

func (s *datingAppService) Swipe(ctx context.Context, input request.SwipeRequest) (*response.SwipeResponse, message.Message) {
	// Validate swipe action.
	if input.Action != "like" && input.Action != "pass" {
		return nil, message.Message{
			Code:    message.ErrReqParam.Code,
			Message: "action must be 'like' or 'pass'",
		}
	}

	// Retrieve the user based on the provided UserID.
	user, err := s.repository.GetUserByID(input.UserID)
	if err != nil {
		s.app.Log.Warn(err.Error())
		return nil, message.ErrUserNotFound
	}

	today := time.Now().Truncate(24 * time.Hour)
	_, err = s.repository.CheckSwipedUser(input.UserID, input.TargetUserID, today)
	if err == nil {
		return nil, message.Message{
			Code:    message.ErrReqParam.Code,
			Message: "already swiped on this profile today",
		}
	}

	// Count how many swipes the user has done today.
	count, err := s.repository.CountSwipedUser(input.UserID, today)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}
	if !user.IsPremium && int(count) >= s.app.Config.AppConfig.SwipeLimit {
		return nil, message.Message{
			Code:    message.ErrReqParam.Code,
			Message: "daily swipe limit reached",
		}
	}
	swipe := entity.Swipe{
		UserID:       input.UserID,
		TargetUserID: input.TargetUserID,
		Action:       input.Action,
		SwipeDate:    today,
	}
	_, err = s.repository.Swipe(swipe)
	if err != nil {
		s.app.Log.Error(err.Error())
		return nil, message.FailedMsg
	}

	return &response.SwipeResponse{
		Message: "Swipe recorded successfully",
	}, message.SuccessMsg
}
