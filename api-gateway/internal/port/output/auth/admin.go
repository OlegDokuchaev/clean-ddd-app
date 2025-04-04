package auth

type Auth interface {
	Validate(token string) bool
}
