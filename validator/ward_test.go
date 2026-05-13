package validators

import (
	"testing"

	ward "github.com/rah-0/ward"
	"github.com/rah-0/ward/types/strs"
)

// --- Single field ---

func BenchmarkWardSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := validUser.Email
		v := ward.New()
		v.Add(strs.New("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
	}
}

func BenchmarkWardSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := invalidUser.Email
		v := ward.New()
		v.Add(strs.New("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
	}
}

func BenchmarkWardMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form := validUser
		v := ward.New()
		v.Add(
			strs.New("Email",    &form.Email,    strs.RuleNotEmpty(), strs.RuleIsEmail()),
			strs.New("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
			strs.New("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
			strs.New("Age",      &form.Age,      strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
			strs.New("Website",  &form.Website,  strs.RuleNotEmpty(), strs.RuleIsURL()),
		)
		v.Run()
	}
}

func BenchmarkWardMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form := invalidUser
		v := ward.New()
		v.Add(
			strs.New("Email",    &form.Email,    strs.RuleNotEmpty(), strs.RuleIsEmail()),
			strs.New("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
			strs.New("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
			strs.New("Age",      &form.Age,      strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
			strs.New("Website",  &form.Website,  strs.RuleNotEmpty(), strs.RuleIsURL()),
		)
		v.Run()
	}
}

func BenchmarkWardMultiFieldStopOnFail(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form := invalidUser
		v := ward.New()
		v.Policy.StopOnFail = true
		v.Add(
			strs.New("Email",    &form.Email,    strs.RuleNotEmpty(), strs.RuleIsEmail()),
			strs.New("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
			strs.New("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
			strs.New("Age",      &form.Age,      strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
			strs.New("Website",  &form.Website,  strs.RuleNotEmpty(), strs.RuleIsURL()),
		)
		v.Run()
	}
}

// --- Full cycle: validate + inspect results ---

func BenchmarkWardFullCycleSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := validUser.Email
		v := ward.New()
		v.Add(strs.New("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
		_ = v.HasFailures()
	}
}

func BenchmarkWardFullCycleSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := invalidUser.Email
		v := ward.New()
		v.Add(strs.New("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				_ = f.FieldName
				_ = f.RuleID
			}
		}
	}
}

func BenchmarkWardFullCycleMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form := validUser
		v := ward.New()
		v.Add(
			strs.New("Email",    &form.Email,    strs.RuleNotEmpty(), strs.RuleIsEmail()),
			strs.New("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
			strs.New("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
			strs.New("Age",      &form.Age,      strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
			strs.New("Website",  &form.Website,  strs.RuleNotEmpty(), strs.RuleIsURL()),
		)
		v.Run()
		_ = v.HasFailures()
	}
}

func BenchmarkWardFullCycleMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form := invalidUser
		v := ward.New()
		v.Add(
			strs.New("Email",    &form.Email,    strs.RuleNotEmpty(), strs.RuleIsEmail()),
			strs.New("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
			strs.New("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
			strs.New("Age",      &form.Age,      strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
			strs.New("Website",  &form.Website,  strs.RuleNotEmpty(), strs.RuleIsURL()),
		)
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				_ = f.FieldName
				_ = f.RuleID
			}
		}
	}
}
