package validators

import (
	"testing"

	ward "github.com/rah-0/ward"
	"github.com/rah-0/ward/types/strs"
)

func wardField(name string, value *string, rules ...strs.Rule) *strs.Field {
	field := strs.New(name, value, rules...)
	// Match the peers: collect the first failing rule for each field.
	field.Policy.StopOnFail = true
	return field
}

func wardAddFields(v *ward.Validate, form *UserForm) {
	v.Add(
		wardField("Email", &form.Email, strs.RuleNotEmpty(), strs.RuleIsEmail()),
		wardField("Username", &form.Username, strs.RuleNotEmpty(), strs.RuleLengthMin(5), strs.RuleLengthMax(50), strs.RuleIsUsernameChars()),
		wardField("Password", &form.Password, strs.RuleNotEmpty(), strs.RuleLengthMin(8), strs.RuleLengthMax(150), strs.RuleIsPasswordChars()),
		wardField("Age", &form.Age, strs.RuleNotEmpty(), strs.RuleIsNonNegativeInt()),
		wardField("Website", &form.Website, strs.RuleNotEmpty(), strs.RuleIsURL()),
	)
}

// --- Single field ---

func BenchmarkWardSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		email := validUser.Email
		v := ward.New()
		v.Add(wardField("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
	}
}

func BenchmarkWardSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		email := invalidUser.Email
		v := ward.New()
		v.Add(wardField("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
	}
}

func BenchmarkWardMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		form := validUser
		v := ward.New()
		wardAddFields(v, &form)
		v.Run()
	}
}

func BenchmarkWardMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		form := invalidUser
		v := ward.New()
		wardAddFields(v, &form)
		v.Run()
	}
}

func BenchmarkWardMultiFieldStopOnFail(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		form := invalidUser
		v := ward.New()
		v.Policy.StopOnFail = true
		wardAddFields(v, &form)
		v.Run()
	}
}

// --- Full cycle: validate + inspect results ---

func BenchmarkWardFullCycleSingleFieldValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		email := validUser.Email
		v := ward.New()
		v.Add(wardField("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
		validationOKSink = !v.HasFailures()
	}
}

func BenchmarkWardFullCycleSingleFieldInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		email := invalidUser.Email
		v := ward.New()
		v.Add(wardField("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail()))
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				validationFieldSink = f.FieldName
				validationRuleIDSink = f.RuleID
			}
		}
	}
}

func BenchmarkWardFullCycleMultiFieldAllValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		form := validUser
		v := ward.New()
		wardAddFields(v, &form)
		v.Run()
		validationOKSink = !v.HasFailures()
	}
}

func BenchmarkWardFullCycleMultiFieldSomeInvalid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		form := invalidUser
		v := ward.New()
		wardAddFields(v, &form)
		v.Run()
		if v.HasFailures() {
			for _, f := range v.Failures() {
				validationFieldSink = f.FieldName
				validationRuleIDSink = f.RuleID
			}
		}
	}
}
