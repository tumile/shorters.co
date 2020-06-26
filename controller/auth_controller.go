package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"shorters/domain"
	"shorters/service/cache"
	"shorters/service/dto"
	"shorters/service/jwt"
	"shorters/service/mail"
	"time"

	"github.com/teris-io/shortid"
)

type AuthController interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	AuthenticateMiddleware(next http.Handler) http.Handler
}

type authController struct {
	jwtService  jwt.JWTService
	mailService mail.MailService
	timeCache   cache.TimeCache
}

func NewAuthController(jwtService jwt.JWTService, mailService mail.MailService) AuthController {
	return &authController{
		jwtService:  jwtService,
		mailService: mailService,
		timeCache:   cache.NewTimeCache(),
	}
}

func (c *authController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Email == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	otp := shortid.MustGenerate()
	err = c.mailService.SendOTP(user.Email, otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.timeCache.Put(user.Email, otp, time.Now().Add(5*time.Minute).Unix())
	w.WriteHeader(http.StatusOK)
}

func (c *authController) Verify(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req.Email == "" || req.OTP == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	otp := c.timeCache.Get(req.Email)
	if otp != req.OTP {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := c.jwtService.Sign(domain.User{Email: req.Email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.TokenResponse{Token: token})
}

func (c *authController) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		user, _ := c.jwtService.Parse(tokenString)
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
