package response

// swagger:model SignUpResponse
type SignUpResponse struct {
	// in: string
	Message string `json:"message,omitempty"`
}

// swagger:model LoginResponse
type LoginResponse struct {
	// in: string
	Token string `json:"token,omitempty"`
}

// swagger:model SwipeResponse
type SwipeResponse struct {
	// in: string
	Message string `json:"message,omitempty"`
}
