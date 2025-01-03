package service

import (
	"os"
	"testing"

	"github.com/emrgen/unpost/internal/tester"
)

func TestMain(m *testing.M) {
	purge, err := tester.SetupDocker()
	if err != nil {
		panic(err)
	}
	defer purge()

	code := m.Run()
	os.Exit(code)
}
