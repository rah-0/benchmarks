package validators

import (
	"slices"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/validator/v10"
	"github.com/rah-0/ward"
	"github.com/rah-0/ward/types/strs"
)

type validationOutcome struct {
	library string
	fields  []string
}

type fieldCase struct {
	name    string
	field   string
	value   string
	invalid bool
}

func checkValidationOutcomes(t *testing.T, input UserForm, want []string) {
	t.Helper()
	var playgroundFields, ozzoFields, wardFields []string
	if err := pgNew(t).Struct(pgLoad(input)); err != nil {
		failures, ok := err.(validator.ValidationErrors)
		if !ok {
			t.Fatalf("unexpected playground error: %v", err)
		}
		for _, failure := range failures {
			playgroundFields = append(playgroundFields, strings.ToLower(failure.Field()))
		}
	}
	if err := ozzoValidateAll(input); err != nil {
		failures, ok := err.(validation.Errors)
		if !ok {
			t.Fatalf("unexpected ozzo error: %v", err)
		}
		for field := range failures {
			ozzoFields = append(ozzoFields, field)
		}
	}
	v := ward.New()
	wardAddFields(v, &input)
	v.Run()
	for _, failure := range v.Failures() {
		wardFields = append(wardFields, strings.ToLower(failure.FieldName))
	}
	for _, outcome := range []validationOutcome{
		{library: "playground", fields: playgroundFields},
		{library: "ozzo", fields: ozzoFields},
		{library: "ward", fields: wardFields},
	} {
		slices.Sort(outcome.fields)
		if !slices.Equal(outcome.fields, want) {
			t.Errorf("%s failure fields: got %v (%d), want %v (%d)", outcome.library, outcome.fields, len(outcome.fields), want, len(want))
		}
	}
}

func TestBenchmarkFixtures(t *testing.T) {
	t.Run("valid", func(t *testing.T) { checkValidationOutcomes(t, validUser, nil) })
	t.Run("invalid", func(t *testing.T) {
		checkValidationOutcomes(t, invalidUser, []string{"age", "email", "password", "username", "website"})
	})
}

func TestComparableFieldRules(t *testing.T) {
	cases := []fieldCase{
		{name: "username minimum", field: "username", value: "abcde"},
		{name: "username maximum", field: "username", value: strings.Repeat("a", 50)},
		{name: "username punctuation", field: "username", value: "a_.@-"},
		{name: "username short", field: "username", value: "abcd", invalid: true},
		{name: "username long", field: "username", value: strings.Repeat("a", 51), invalid: true},
		{name: "username space", field: "username", value: "john doe", invalid: true},
		{name: "username unicode", field: "username", value: "johné", invalid: true},
		{name: "password minimum", field: "password", value: "Aa1@xxxx"},
		{name: "password maximum", field: "password", value: "Aa1@" + strings.Repeat("x", 146)},
		{name: "password rune length", field: "password", value: "Aa1@" + strings.Repeat("é", 146)},
		{name: "password short", field: "password", value: "Aa1@xxx", invalid: true},
		{name: "password long", field: "password", value: "Aa1@" + strings.Repeat("x", 147), invalid: true},
		{name: "password missing lowercase", field: "password", value: "AA1@XXXX", invalid: true},
		{name: "password missing uppercase", field: "password", value: "aa1@xxxx", invalid: true},
		{name: "password missing digit", field: "password", value: "Aaa@xxxx", invalid: true},
		{name: "password missing special", field: "password", value: "Aa11xxxx", invalid: true},
		{name: "age zero", field: "age", value: "0"},
		{name: "age maximum", field: "age", value: "999999999"},
		{name: "age too large", field: "age", value: "1000000000", invalid: true},
		{name: "age negative", field: "age", value: "-5", invalid: true},
		{name: "age leading zero", field: "age", value: "01", invalid: true},
		{name: "age positive sign", field: "age", value: "+5", invalid: true},
		{name: "age fractional", field: "age", value: "1.5", invalid: true},
	}
	for _, field := range []string{"email", "username", "password", "age", "website"} {
		cases = append(cases, fieldCase{name: field + " required", field: field, invalid: true})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := validUser
			switch tc.field {
			case "email":
				input.Email = tc.value
			case "username":
				input.Username = tc.value
			case "password":
				input.Password = tc.value
			case "age":
				input.Age = tc.value
			case "website":
				input.Website = tc.value
			}
			var want []string
			if tc.invalid {
				want = []string{tc.field}
			}
			checkValidationOutcomes(t, input, want)
		})
	}
}

func TestSingleEmailAndStopOnFail(t *testing.T) {
	for _, email := range []string{validUser.Email, invalidUser.Email} {
		want := email == invalidUser.Email
		if got := pgNew(t).Struct(pgEmailForm{Email: email}) != nil; got != want {
			t.Errorf("playground email failure: got %t, want %t", got, want)
		}
		if got := ozzoValidateEmail(email) != nil; got != want {
			t.Errorf("ozzo email failure: got %t, want %t", got, want)
		}
		v := ward.New().Add(wardField("Email", &email, strs.RuleNotEmpty(), strs.RuleIsEmail())).Run()
		if v.HasFailures() != want || len(v.Failures()) > 1 {
			t.Errorf("ward email failures: got %v, want failure=%t", v.Failures(), want)
		}
	}
	input := invalidUser
	v := ward.New()
	v.Policy.StopOnFail = true
	wardAddFields(v, &input)
	v.Run()
	if failures := v.Failures(); len(failures) != 1 || failures[0].FieldName != "Email" {
		t.Fatalf("stop-on-fail must report only Email, got %v", failures)
	}
}
