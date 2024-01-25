package secrets

func getTestSecrets() Secrets {
	return &testSecrets{}
}

type testSecrets struct {
}

func (s *testSecrets) Get(_ string) (string, error) {
	return "test-secret-value", nil
}
