package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type pgUserForm struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=5,max=50,alphanumunicode"`
	Password string `validate:"required,min=8,max=150"`
	Age      string `validate:"required,numeric"`
	Website  string `validate:"required,url"`
}

type pgEmailForm struct {
	Email string `validate:"required,email"`
}

func pgLoad(u UserForm) pgUserForm {
	return pgUserForm{Email: u.Email, Username: u.Username, Password: u.Password, Age: u.Age, Website: u.Website}
}

func BenchmarkPlaygroundSingleFieldValid(b *testing.B) {
	v := validator.New()
	s := pgEmailForm{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Email = validUser.Email
		_ = v.Struct(s)
	}
}

func BenchmarkPlaygroundSingleFieldInvalid(b *testing.B) {
	v := validator.New()
	s := pgEmailForm{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Email = invalidUser.Email
		_ = v.Struct(s)
	}
}

func BenchmarkPlaygroundMultiFieldAllValid(b *testing.B) {
	v := validator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(pgLoad(validUser))
	}
}

func BenchmarkPlaygroundMultiFieldSomeInvalid(b *testing.B) {
	v := validator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(pgLoad(invalidUser))
	}
}
