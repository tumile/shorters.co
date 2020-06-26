package dto

type VerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type StoreRequest struct {
	URL string `json:"url"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type StoreResponse struct {
	Key string `json:"key"`
}

type LinkResponse struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	Visits int64  `json:"visits"`
}
