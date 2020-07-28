package dto

type VerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type CustomShortenRequest struct {
	URL       string `json:"url"`
	CustomKey string `json:"custom_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ShortenResponse struct {
	Key string `json:"key"`
}

type LinkResponse struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	Visits int64  `json:"visits"`
}
