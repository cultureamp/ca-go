// Package evaluationcontext defines the attributes for contexts you can
// provide in a query for a flag or product toggle. Constructor functions are
// exposed to create valid instances of evaluation contexts.
//
// An evaluation context is simply a bag of attributes keyed by a unique
// identifier. These values are used in two places:
//  1. When creating a flag or segment, the attributes are used to form the
//     targeting rules. For example, "if the user's realUserID is 123, return
//     false for this flag".
//  2. When querying for a flag in your service, the SDK uses the attributes
//     to evaluate the rules for the flag to return the correct value. Using
//     the same example as above, supplying an EvaluationContext with a realUserID
//     of 123 would cause the flag to evaluate to false.
//
// You should always supply as many attributes as you can to give yourself more
// flexibility when writing new targeting rules. When you query a flag containing a
// rule that works on attribute "foo", you must supply attribute "foo" in the
// evaluation context.
//
// User and Survey are now deprecated but have not been removed for backwards
// compatibility. In order to upgrade use EvaluationContext instead there are
// three steps to follow:
//  1. Install the latest version of this package.
//  2. Update targeting rules in LD to use contexts as explained on [confluence].
//  3. Change all uses of User and Survey helper functions to EvaluationContext.
//
// [confluence]: https://cultureamp.atlassian.net/wiki/spaces/CoSe/pages/2820768174/Guide+to+flag+targeting+in+LaunchDarkly
package evaluationcontext
