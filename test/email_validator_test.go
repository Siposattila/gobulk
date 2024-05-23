package email_validator_test

import (
	"testing"

	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
)

func TestValidatorOnGoodEmail(t *testing.T) {
	e := email.Email{Email: "sattipolo@gmail.com"}
	e.ValidateEmail()

	if e.Valid == interfaces.EMAIL_INVALID {
		t.Fatalf("e.Valid should be true on sattipolo@gmail.com")
	}
}

func TestValidatorOnBadEmail(t *testing.T) {
	e := email.Email{Email: "iepwhfepwifhewipfeh@gmail.com"}
	e.ValidateEmail()

	if e.Valid == interfaces.EMAIL_VALID {
		t.Fatalf("e.Valid should not be true on iepwhfepwifhewipfeh@gmail.com")
	}
}

func TestValidatorOnGmail(t *testing.T) {
	e := email.Email{Email: "gobulk189@gmail.com"}
	e.ValidateEmail()

	if e.Valid == interfaces.EMAIL_INVALID {
		t.Fatalf("e.Valid should be true on gobulk189@gmail.com")
	}
}

func TestValidatorOnFreemail(t *testing.T) {
	e := email.Email{Email: "gobulk@fremail.hu"}
	e.ValidateEmail()

	if e.Valid == interfaces.EMAIL_INVALID {
		t.Fatalf("e.Valid should be true on gobulk@fremail.hu")
	}
}
