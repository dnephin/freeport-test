package one

import (
	"testing"

	freeporttest "github.com/dnephin/freeport-test"
)

func TestConflicts(t *testing.T) {
	freeporttest.RunTestConflicts(t, "once", 1)
}
