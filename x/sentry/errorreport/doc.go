// Package errorreport enables you to configure Sentry for error reporting. A
// general ReportError function is provided for ad-hoc reporting, as well as
// several options for middleware. These middleware detect and report errors
// automatically.
//
// See the executable examples below and on individual functions for more
// details on usage.
//
// You configure and initialise errorreport using Init():
//   err := errorreport.Init(
//            errorreport.WithDSN(os.Getenv("SENTRY_DSN")),
//            errorreport.WithRelease(os.Getenv("APP"), os.Getenv("APP_VERSION")),
//            errorreport.WithEnvironment(os.Getenv("AWS_ENVIRONMENT_NAME")))
//   if err != nil {
//     // handle initialisation error
//   }
//
// Ad-hoc errors can be reported using ReportError():
//   errorreport.ReportError(ctx, errors.New("We hit a snag!"))
//
// HTTP middleware can be used. Passing in nil uses the default panic handler.
// See the OnRequestPanicHandler type if you wish to supply your own.
//   mw := middleware.NewHTTPMiddleware(nil)
//   mw(myHTTPHandler)
//
// Goa middleware can be used, allowing errors to be reported automatically from
// Goa applications:
//   mw := errorreport.NewGoaMiddleware()
//   mw(myGoaEndpoint)
//
// This is recommended when using Goa, as it offers reporting of all errors
// returned from the generated logic types.
//
// This package also supports middleware for Lambda functions:
//   mw := errorreport.NewLambdaMiddleware(errorreport.LambdaErrorOptions{})
package errorreport
