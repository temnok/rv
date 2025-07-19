package rv

import (
	"testing"
)

func TestRunKernelDirect64(t *testing.T) {
	testRunKernel(t, 64, "linux/output", "user@rv64")
}
