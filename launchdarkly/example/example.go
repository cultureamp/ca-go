package main

import (
	"fmt"

	flags "github.com/cultureamp/ca-go/launchdarkly"
	"github.com/cultureamp/ca-go/launchdarkly/evaluationcontext"
)

func main() {
	contexts := map[string]evaluationcontext.Context{
		"anonymous":  evaluationcontext.NewEvaluationContext(),
		"user":       evaluationcontext.NewEvaluationContext(evaluationcontext.WithUserID("go-user-1"), evaluationcontext.WithContextRealUserID("go-real-user-1")),
		"survey":     evaluationcontext.NewEvaluationContext(evaluationcontext.WithSurveyID("go-survey-2")),
		"account":    evaluationcontext.NewEvaluationContext(evaluationcontext.WithAccountID("go-account-3")),
		"subdomain":  evaluationcontext.NewAnonymousContextWithSubdomain("go-account-10", "goSubdomain"),
		"everything": evaluationcontext.NewEvaluationContext(evaluationcontext.WithUserID("go-user-4"), evaluationcontext.WithContextRealUserID("go-real-user-4"), evaluationcontext.WithAccountID("go-account-4"), evaluationcontext.WithSurveyID("go-survey-4")),
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
