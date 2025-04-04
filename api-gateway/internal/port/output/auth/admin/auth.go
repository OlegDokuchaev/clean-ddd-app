package admin

type Auth interface {
	Validate(token string) bool
}
