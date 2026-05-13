package validators

// UserForm is the common input struct used across all validator benchmarks.
// All fields are strings to keep the comparison apples-to-apples.
type UserForm struct {
	Email    string
	Username string
	Password string
	Age      string
	Website  string
}

var validUser = UserForm{
	Email:    "user@example.com",
	Username: "john_doe",
	Password: "Str0ng@Pass",
	Age:      "25",
	Website:  "https://example.com",
}

var invalidUser = UserForm{
	Email:    "not-an-email",
	Username: "jo",
	Password: "weak",
	Age:      "-5",
	Website:  "not-a-url",
}
