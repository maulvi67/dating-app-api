package message

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var SuccessMsg = Message{Code: 200, Message: "Success"}
var FailedMsg = Message{Code: 500, Message: "Failed"}
var ErrDataExists = Message{Code: 400, Message: "Data already exists"}
var ErrReqParam = Message{Code: 400, Message: "Param invalid"}
var ErrUserNotFound = Message{Code: 400, Message: "User not found"}
var ErrInvalidCred = Message{Code: 401, Message: "Invalid credentials"}
var AuthenticationFailed = Message{Code: 401, Message: "JWT token is invalid"}
var UnauthorizedError = Message{Code: 401, Message: "No authorization token was found"}
