package methods

type Method interface {
	GetNextServer() string
	AddServer(server string)
	RemoveServer(server string)
	GetServers() []string
	GetServerCount() int
	GetServerHealth() map[string]bool
	SetServerHealth(server string, health bool)
}
