package domain

import "testing"

// Discovery returns small metadata. It must not allocate proportional to every
// embedded request/response example each time a CLI command tree is constructed.
func TestEndpointDiscoveryAllocationsStayBounded(t *testing.T) {
	_ = AllEndpoints()
	allocations := testing.AllocsPerRun(2, func() { _ = AllEndpoints() })
	if allocations > 5000 {
		t.Fatalf("discovery allocated %.0f objects; large source documents should not be reparsed for each invocation", allocations)
	}
}
func TestDiscoveredDocumentationIsIndependent(t *testing.T) {
	e, _ := EndpointByID("cloud-auth", "register")
	e.Documentation.Title = "changed"
	again, _ := EndpointByID("cloud-auth", "register")
	if again.Documentation.Title == "changed" {
		t.Fatal("a caller mutated another discovery result")
	}
}
