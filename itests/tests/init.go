package tests

import "github.com/pegondo/starwars/service/ex/client"

const (
	// addr is the address to target.
	addr = "http://localhost:8080"
	// pageSizeStr is the standard page size for the resources to target.
	pageSize = 15
	// pageSizeStr is the standard page size for the resources to target in
	// string format.
	pageSizeStr = "15"
)

// c is the client to request the target.
var c = client.New(addr)
