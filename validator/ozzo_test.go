package validators

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ozzoValidateEmail(email string) error {
	return validation.Validate(email, validation.Required, is.Email)
}

var passwordRule = validation.By(func(value any) error {
	if passwordChars(value.(string)) {
		return nil
	}
	return validation.NewError("password", "must contain a lowercase letter, uppercase letter, digit and special character")
})

func ozzoValidateAll(u UserForm) error {
	return validation.Errors{
		"email":    validation.Validate(u.Email, validation.Required, is.Email),
		"username": validation.Validate(u.Username, validation.Required, validation.RuneLength(5, 50), validation.Match(usernameChars)),
		"password": validation.Validate(u.Password, validation.Required, validation.RuneLength(8, 150), passwordRule),
		"age":      validation.Validate(u.Age, validation.Required, validation.Match(nonNegativeInt)),
		"website":  validation.Validate(u.Website, validation.Required, is.URL),
	}.Filter()
}

func BenchmarkOzzoSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = ozzoValidateEmail(validUser.Email)
	}
}

func BenchmarkOzzoSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = ozzoValidateEmail(invalidUser.Email)
	}
}

func BenchmarkOzzoMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = ozzoValidateAll(validUser)
	}
}

func BenchmarkOzzoMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = ozzoValidateAll(invalidUser)
	}
}

// --- Full cycle: validate + inspect results ---

func BenchmarkOzzoFullCycleSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationOKSink = ozzoValidateEmail(validUser.Email) == nil
	}
}

func BenchmarkOzzoFullCycleSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if err := ozzoValidateEmail(invalidUser.Email); err != nil {
			validationRuleSink = err.Error()
		}
	}
}

func BenchmarkOzzoFullCycleMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		validationOKSink = ozzoValidateAll(validUser) == nil
	}
}

func BenchmarkOzzoFullCycleMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if err := ozzoValidateAll(invalidUser); err != nil {
			for field, e := range err.(validation.Errors) {
				validationFieldSink = field
				validationRuleSink = e.Error()
			}
		}
	}
}
