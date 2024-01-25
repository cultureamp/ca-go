package secrets

type Secrets interface {
	Get(name string) (string, error)
}
