package jwt

import (
	"reflect"
	"shorters/domain"
	"testing"
)

func TestJWTService(t *testing.T) {
	t.Run("Test JWT Parse", func(t *testing.T) {
		j := NewJWTService()
		user := domain.User{Email: "random@email.com"}
		ss, err := j.Sign(user)
		if err != nil {
			t.Errorf("Sign() error = %v", err)
			return
		}
		userParsed, err := j.Parse(ss)
		if err != nil {
			t.Errorf("Parse() error = %v", err)
			return
		}
		if !reflect.DeepEqual(user, userParsed) {
			t.Errorf("Parse() got = %v, want %v", userParsed, user)
		}
	})
}
