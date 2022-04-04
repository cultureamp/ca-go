package evaluationcontext

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

// Context represents a set of attributes which a flag is evaluated against. The
// only context supported now is User.
type Context interface {
	// ToLDUser transforms the context implementation into an LDUser object that can
	// be understood by LaunchDarkly when evaluating a flag.
	ToLDUser() lduser.User
}
