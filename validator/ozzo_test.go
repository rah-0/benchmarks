package validators

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ozzoValidateEmail(email string) error {
	return validation.Validate(email, validation.Required, is.Email)
}

func ozzoValidateAll(u UserForm) error {
	return validation.Errors{
		"email":    validation.Validate(u.Email, validation.Required, is.Email),
		"username": validation.Validate(u.Username, validation.Required, validation.Length(5, 50)),
		"password": validation.Validate(u.Password, validation.Required, validation.Length(8, 150)),
		"age":      validation.Validate(u.Age, validation.Required, is.Int),
		"website":  validation.Validate(u.Website, validation.Required, is.URL),
	}.Filter()
}

func BenchmarkOzzoSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ozzoValidateEmail(validUser.Email)
	}
}

func BenchmarkOzzoSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ozzoValidateEmail(invalidUser.Email)
	}
}

func BenchmarkOzzoMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ozzoValidateAll(validUser)
	}
}

func BenchmarkOzzoMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ozzoValidateAll(invalidUser)
	}
}
