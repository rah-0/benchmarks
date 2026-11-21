package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type pgUserForm struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=5,max=50,username"`
	Password string `validate:"required,min=8,max=150,password"`
	Age      string `validate:"required,nonnegativeint"`
	Website  string `validate:"required,url"`
}

type pgEmailForm struct {
	Email string `validate:"required,email"`
}

func pgLoad(u UserForm) pgUserForm {
	return pgUserForm{Email: u.Email, Username: u.Username, Password: u.Password, Age: u.Age, Website: u.Website}
}

func pgNew(tb testing.TB) *validator.Validate {
	tb.Helper()
	v := validator.New()
	if err := v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameChars.MatchString(fl.Field().String())
	}); err != nil {
		tb.Fatal(err)
	}
	if err := v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return passwordChars(fl.Field().String())
	}); err != nil {
		tb.Fatal(err)
	}
	if err := v.RegisterValidation("nonnegativeint", func(fl validator.FieldLevel) bool {
		return nonNegativeInt.MatchString(fl.Field().String())
	}); err != nil {
		tb.Fatal(err)
	}
	return v
}

func BenchmarkPlaygroundSingleFieldValid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		s := pgEmailForm{Email: validUser.Email}
		validationErrorSink = v.Struct(s)
	}
}

func BenchmarkPlaygroundSingleFieldInvalid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		s := pgEmailForm{Email: invalidUser.Email}
		validationErrorSink = v.Struct(s)
	}
}

func BenchmarkPlaygroundMultiFieldAllValid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = v.Struct(pgLoad(validUser))
	}
}

func BenchmarkPlaygroundMultiFieldSomeInvalid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		validationErrorSink = v.Struct(pgLoad(invalidUser))
	}
}

// --- Full cycle: validate + inspect results ---

func BenchmarkPlaygroundFullCycleSingleFieldValid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		s := pgEmailForm{Email: validUser.Email}
		validationOKSink = v.Struct(s) == nil
	}
}

func BenchmarkPlaygroundFullCycleSingleFieldInvalid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		s := pgEmailForm{Email: invalidUser.Email}
		if err := v.Struct(s); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				validationFieldSink = e.Field()
				validationRuleSink = e.Tag()
			}
		}
	}
}

func BenchmarkPlaygroundFullCycleMultiFieldAllValid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		validationOKSink = v.Struct(pgLoad(validUser)) == nil
	}
}

func BenchmarkPlaygroundFullCycleMultiFieldSomeInvalid(b *testing.B) {
	v := pgNew(b)
	b.ReportAllocs()
	for b.Loop() {
		if err := v.Struct(pgLoad(invalidUser)); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				validationFieldSink = e.Field()
				validationRuleSink = e.Tag()
			}
		}
	}
}
