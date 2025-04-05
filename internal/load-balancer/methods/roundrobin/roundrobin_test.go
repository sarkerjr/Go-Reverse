package roundrobin

import (
	"testing"
)

func TestNewRoundRobin(t *testing.T) {
	rr := NewRoundRobin()
	if rr == nil {
		t.Fatal("Expected non-nil RoundRobin instance")
	}
	if rr.GetServerCount() != 0 {
		t.Errorf("Expected server count 0, got %d", rr.GetServerCount())
	}
}

func TestAddServer(t *testing.T) {
	rr := NewRoundRobin()
	rr.AddServer("server1")
	if rr.GetServerCount() != 1 {
		t.Errorf("Expected server count 1, got %d", rr.GetServerCount())
	}
	if rr.GetServers()[0] != "server1" {
		t.Errorf("Expected 'server1', got %s", rr.GetServers()[0])
	}
	if !rr.GetServerHealth()["server1"] {
		t.Errorf("Expected server1 to be healthy")
	}
}

func TestRemoveServer(t *testing.T) {
	rr := NewRoundRobin()
	rr.AddServer("server1")
	rr.AddServer("server2")
	rr.RemoveServer("server1")
	if rr.GetServerCount() != 1 {
		t.Errorf("Expected server count 1 after removal, got %d", rr.GetServerCount())
	}
	if rr.GetServers()[0] != "server2" {
		t.Errorf("Expected remaining server to be server2, got %s", rr.GetServers()[0])
	}
	if _, exists := rr.GetServerHealth()["server1"]; exists {
		t.Errorf("Expected server1 to be removed from health map")
	}
}

func TestGetNextServer(t *testing.T) {
	rr := NewRoundRobin()
	servers := []string{"server1", "server2", "server3"}
	for _, s := range servers {
		rr.AddServer(s)
	}

	for i := 0; i < 6; i++ {
		expected := servers[i%len(servers)]
		got := rr.GetNextServer()
		if got != expected {
			t.Errorf("Expected %s, got %s", expected, got)
		}
	}
}

func TestGetNextServerEmpty(t *testing.T) {
	rr := NewRoundRobin()
	if server := rr.GetNextServer(); server != "" {
		t.Errorf("Expected empty string for no servers, got %s", server)
	}
}

func TestSetServerHealth(t *testing.T) {
	rr := NewRoundRobin()
	rr.AddServer("server1")
	rr.SetServerHealth("server1", false)
	if rr.GetServerHealth()["server1"] {
		t.Errorf("Expected server1 to be unhealthy")
	}
}
