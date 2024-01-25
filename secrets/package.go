package secrets

import (
	"flag"
	"os"
	"strings"
)

var defaultSecrets = getInstance()

func getInstance() Secrets {
	if isTestMode() {
		return getTestSecrets()
	}
	region := os.Getenv("AWS_REGION")
	return NewAWSSecrets(region)
}
func Get(name string) (string, error) {
	return defaultSecrets.Get(name)
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}

func SetImpl(impl Secrets) {
	defaultSecrets = impl
}
