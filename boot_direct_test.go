package rv

import (
	"testing"
)

func TestRunKernelDirect64(t *testing.T) {
	testRunKernel(t, 64, "biko/output-64", "user@rv")
}
