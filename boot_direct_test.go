package rv

import (
	"testing"
)

func TestRunKernelDirect64(t *testing.T) {
	testRunKernel(t, 64, "linux/output", "No filesystem could mount root, tried:")
}
