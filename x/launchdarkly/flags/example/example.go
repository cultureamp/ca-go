package main

import (
	"fmt"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags"
	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
)

func main() {
	// skip staticcheck linter for this block as there are deprecated methods used, can be removed when legacy code is removed
	//nolint:staticcheck
	contexts := map[string]evaluationcontext.Context{
		"anonymous":  evaluationcontext.NewEvaluationContext(),
		"user":       evaluationcontext.NewEvaluationContext(evaluationcontext.WithUserID("go-user-1"), evaluationcontext.WithContextRealUserID("go-real-user-1")),
		"survey":     evaluationcontext.NewEvaluationContext(evaluationcontext.WithSurveyID("go-survey-2")),
		"account":    evaluationcontext.NewEvaluationContext(evaluationcontext.WithAccountID("go-account-3")),
		"subdomain":  evaluationcontext.NewAnonymousContextWithSubdomain("go-account-10", "goSubdomain"),
		"everything": evaluationcontext.NewEvaluationContext(evaluationcontext.WithUserID("go-user-4"), evaluationcontext.WithContextRealUserID("go-real-user-4"), evaluationcontext.WithAccountID("go-account-4"), evaluationcontext.WithSurveyID("go-survey-4")),
		// todo: remove these examples once deprecated legacy code is deleted
		"legacy user": evaluationcontext.NewUser("go-user-5", evaluationcontext.WithRealUserID("go-real-user-5")),
		"legacy user with account": evaluationcontext.NewUser("go-user-6", evaluationcontext.WithRealUserID("go-real-user-6"),
			evaluationcontext.WithUserAccountID("go-account-6")),
		"legacy anonymous user":          evaluationcontext.NewAnonymousUser(""),
		"legacy anonymous user with key": evaluationcontext.NewAnonymousUser("go-user-7"),
		"legacy survey":                  evaluationcontext.NewSurvey("go-survey-8"),
		"legacy survey with account":     evaluationcontext.NewSurvey("go-survey-9", evaluationcontext.WithSurveyAccountID("go-account-9")),
	}

	err := flags.Configure()
	if err != nil {
		fmt.Println(err)
	}
	client, err := flags.GetDefaultClient()
	if err != nil {
		fmt.Println(err)
	}
	defer func(client *flags.Client) {
		err := client.Shutdown()
		if err != nil {
			fmt.Printf("error shutting down client %e", err)
		}
	}(client)
	for key, element := range contexts {
		value, err := client.QueryBoolWithEvaluationContext("test.mode.flag", element, false)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Context: %v\n Flag value: %t\n", key, value)
	}
}
