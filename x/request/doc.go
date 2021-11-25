// Package request exposes types and helper methods to create, add, and retrieve
// request-scoped attributes to context.Context.
//
// Request-scoped attributes include identifiers like the request and correlation
// IDs. When the request is authenticated, user identifiers like the account and
// user aggregate IDs can also be added to the context.
package request
