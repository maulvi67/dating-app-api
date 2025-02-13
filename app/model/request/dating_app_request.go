package request

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SwipeRequest struct {
	TargetUserID uint   `json:"target_user_id"`
	Action       string `json:"action"`
	// userID is added from the JWT middleware.
	UserID uint `json:"-"`
}

type PurchasePremiumRequest struct {
	PackageType string `json:"package_type"`
	UserID      uint   `json:"-"`
}

// swagger:parameters ReqSignUpRequestBody
type ReqSignUpRequestBody struct {
	//  in: body
	Body SignUpRequest `json:"body"`
}

// swagger:parameters ReqLoginRequestBody
type ReqLoginRequestBody struct {
	//  in: body
	Body LoginRequest `json:"body"`
}

// swagger:parameters ReqSwipeRequestBody
type ReqSwipeRequestBody struct {
	//  in: body
	Body SwipeRequest `json:"body"`
}
