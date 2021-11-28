package tre

import (
	"testing"

	freeporttest "github.com/dnephin/freeport-test"
)

func TestConflicts(t *testing.T) {
	freeporttest.RunTestConflicts(t, "thre", 7)
}
