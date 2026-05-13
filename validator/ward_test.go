package validators

import (
	"testing"

	"github.com/rah-0/ward/types/strs"
	wardvalidator "github.com/rah-0/ward/validator"
)

var (
	form         UserForm
	wardEmail    = strs.New("Email",    &form.Email,    strs.NotEmpty(), strs.IsEmail())
	wardUsername = strs.New("Username", &form.Username, strs.NotEmpty(), strs.LengthMin(5), strs.LengthMax(50), strs.IsUsernameChars())
	wardPassword = strs.New("Password", &form.Password, strs.NotEmpty(), strs.LengthMin(8), strs.LengthMax(150), strs.IsPasswordChars())
	wardAge      = strs.New("Age",      &form.Age,      strs.NotEmpty(), strs.IsNonNegativeInt())
	wardWebsite  = strs.New("Website",  &form.Website,  strs.NotEmpty(), strs.IsURL())
)

// --- Single field ---

func BenchmarkWardSingleFieldValid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.Email = validUser.Email
		v.Add(wardEmail)
		v.Run()
		v.Reset()
	}
}

func BenchmarkWardSingleFieldInvalid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.Email = invalidUser.Email
		v.Add(wardEmail)
		v.Run()
		v.Reset()
	}
}

func BenchmarkWardMultiFieldAllValid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form = validUser
		v.Add(wardEmail)
		v.Add(wardUsername)
		v.Add(wardPassword)
		v.Add(wardAge)
		v.Add(wardWebsite)
		v.Run()
		v.Reset()
	}
}

func BenchmarkWardMultiFieldSomeInvalid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form = invalidUser
		v.Add(wardEmail)
		v.Add(wardUsername)
		v.Add(wardPassword)
		v.Add(wardAge)
		v.Add(wardWebsite)
		v.Run()
		v.Reset()
	}
}

func BenchmarkWardMultiFieldStopOnFail(b *testing.B) {
	v := wardvalidator.New()
	v.Policy.StopOnFail = true
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form = invalidUser
		v.Add(wardEmail)
		v.Add(wardUsername)
		v.Add(wardPassword)
		v.Add(wardAge)
		v.Add(wardWebsite)
		v.Run()
		v.Reset()
	}
}

// --- Full cycle: validate + inspect results ---

func BenchmarkWardFullCycleSingleFieldValid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.Email = validUser.Email
		v.Add(wardEmail)
		v.Run()
		_ = v.HasFailures()
		v.Reset()
	}
}

func BenchmarkWardFullCycleSingleFieldInvalid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.Email = invalidUser.Email
		v.Add(wardEmail)
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				_ = f.FieldName
				_ = f.RuleID
			}
		}
		v.Reset()
	}
}

func BenchmarkWardFullCycleMultiFieldAllValid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form = validUser
		v.Add(wardEmail)
		v.Add(wardUsername)
		v.Add(wardPassword)
		v.Add(wardAge)
		v.Add(wardWebsite)
		v.Run()
		_ = v.HasFailures()
		v.Reset()
	}
}

func BenchmarkWardFullCycleMultiFieldSomeInvalid(b *testing.B) {
	v := wardvalidator.New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form = invalidUser
		v.Add(wardEmail)
		v.Add(wardUsername)
		v.Add(wardPassword)
		v.Add(wardAge)
		v.Add(wardWebsite)
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				_ = f.FieldName
				_ = f.RuleID
			}
		}
		v.Reset()
	}
}
