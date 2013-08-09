package ftlog

import (
	"regexp"
)

const (
	// Take care: this pattern is used in xlmfilter, possibly elsewhere
	// This is RFC2822's atext.

	ATEXT   = `[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+`
	AT_FREE = ATEXT + `(?:.` + ATEXT + `)*`

	// This permits an RFC2822 message ID but is a little less restrictive
	PATH_PAT = AT_FREE + `(?:@` + AT_FREE + `)?`

	BODY_LINE_1_PAT = `^(\d+) ([0-9a-fA-F]{40}) ([0-9a-fA-F]{40}) "([^"]*)" (` +
		PATH_PAT + `)$`

	BODY_LINE_3_PAT = `^(\d+) ([0-9a-fA-F]{64}) ([0-9a-fA-F]{64}) "([^"]*)" (` +
		PATH_PAT + `)$`

	IGNORABLE_PAT = `(^ *$)|^ *#`
)

var (
	ignorableRE = regexp.MustCompile(IGNORABLE_PAT)
	pathRE      = regexp.MustCompile(PATH_PAT)
	bodyLine1RE = regexp.MustCompile(BODY_LINE_1_PAT)
	bodyLine3RE = regexp.MustCompile(BODY_LINE_3_PAT)
)

func IgnorableRE() *regexp.Regexp {
	return ignorableRE
}

// This is not necessarily a POSIX path.  In fact it permits most
// or all email addresses as well.
func PathRE() *regexp.Regexp {
	return pathRE
}

func BodyLine1RE() *regexp.Regexp {
	return bodyLine1RE
}

func BodyLine3RE() *regexp.Regexp {
	return bodyLine3RE
}
