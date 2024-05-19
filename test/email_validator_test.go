package email_validator_test

import (
	"testing"

	"github.com/Siposattila/gobulk/internal/email"
)

func TestValidatorOnGoodEmail(t *testing.T) {
	e := email.Email{Email: "sattipolo@gmail.com"}
	e.ValidateEmail()

	if e.Valid == email.EMAIL_INVALID {
		t.Fatalf("e.Valid should be true on sattipolo@gmail.com")
	}
}

func TestValidatorOnBadEmail(t *testing.T) {
	e := email.Email{Email: "iepwhfepwifhewipfeh@gmail.com"}
	e.ValidateEmail()

	if e.Valid == email.EMAIL_VALID {
		t.Fatalf("e.Valid should not be true on iepwhfepwifhewipfeh@gmail.com")
	}
}
