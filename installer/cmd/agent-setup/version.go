package main

// version is set at build time via:
//
//	go build -ldflags "-X main.version=vX.Y.Z" ./cmd/agent-setup
//
// Dev builds default to "dev" and skip the update check.
var version = "dev"
