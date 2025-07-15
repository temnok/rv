package rv

import (
	"testing"
)

func xTestRunKernelDirect64(t *testing.T) {
	testRunKernel(t, 64, "linux/output", "Run /init as init process")
}
