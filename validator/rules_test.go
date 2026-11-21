package validators

import "regexp"

// These predicates match the username, password and integer rules used by Ward.
// Email and URL syntax are checked by each library's native validator.
var (
	usernameChars        = regexp.MustCompile(`^[a-zA-Z0-9_.@-]+$`)
	nonNegativeInt       = regexp.MustCompile(`^(0|[1-9]\d{0,8})$`)
	lowercase            = regexp.MustCompile(`[a-z]`)
	uppercase            = regexp.MustCompile(`[A-Z]`)
	digit                = regexp.MustCompile(`\d`)
	specialChar          = regexp.MustCompile(`[@$!%*?&_\-=]`)
	validationErrorSink  error
	validationOKSink     bool
	validationFieldSink  string
	validationRuleSink   string
	validationRuleIDSink uint32
)

func passwordChars(value string) bool {
	return lowercase.MatchString(value) && uppercase.MatchString(value) &&
		digit.MatchString(value) && specialChar.MatchString(value)
}
