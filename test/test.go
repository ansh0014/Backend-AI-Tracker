package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code here, if needed
	// For example, you can initialize a test database or mock services

	// Run the tests
	exitCode := m.Run()

	// Teardown code here, if needed
	// For example, you can close database connections or clean up resources

	// Exit with the appropriate code
	os.Exit(exitCode)
}
