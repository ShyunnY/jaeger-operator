package logging

import "testing"

func TestLogger(t *testing.T) {

	logger := NewLogger("debug")
	logger.WithName("reconcile").Info("this is reconcile")
	logger.WithName("status").Info("this is status")
	logger.WithValues("runner", "infra").Info("this is infra")
}
