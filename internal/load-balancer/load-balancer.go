package loadbalancer

type ILoadBalancer interface {
	GetNextServer() string
	AddServer(server string)
	RemoveServer(server string)
	GetServers() []string
	GetServerCount() int
	GetServerStatus() map[string]string
	GetServerHealth() map[string]bool
	SetServerHealth(server string, health bool)
}
