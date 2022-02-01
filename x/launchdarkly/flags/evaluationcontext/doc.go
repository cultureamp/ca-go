// Package evaluationcontext defines the kinds of evaluation contexts you can
// provide in a query for a flag or product toggle. Constructor functions are
// exposed to create valid instances of evaluation contexts.
//
// An evaluation context is simply a bag of attributes keyed by a unique
// identifier. These values are used in two situations:
//   1. When creating a flag or segment, the attributes are used to form the
//      targeting rules. For example, "if the user's accountID is 123, return
//      false for this flag".
//   2. When querying for a flag in your service, the SDK uses the attributes
//      to evaluate the rules for the flag to return the correct value. Using
//      the same example as above, supplying a User context with an accountID
//      of 123 would cause the flag to evaluate to false.
//
// Some evaluation contexts (like the User) have constructor functions which
// allow you to supply optional attributes. You should always supply as many
// attributes as you can to give yourself more flexibility when writing new
// targeting rules. When you query a flag containing a rule that works on
// attribute "foo", you must supply attribute "foo" in the evaluation context.
//
// The constructor functions will namespace the key of the evaluation context
// you send in the query. You do not need to prefix the values provided to
// the constructor functions with the entity type. For example, supply the user
// ID as-is rather than prefixing as `user.<user-id>`.
package evaluationcontext
