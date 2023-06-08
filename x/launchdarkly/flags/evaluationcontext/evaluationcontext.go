package evaluationcontext

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

// Context represents a set of attributes which a flag is evaluated against. The
// only contexts supported now are User and Survey
type Context interface {
	// ToLDContext transforms the context implementation into an LDContext object that can
	// be understood by LaunchDarkly when evaluating a flag.
	ToLDContext() ldcontext.Context
}
