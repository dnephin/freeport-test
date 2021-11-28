package freeporttest

import (
	"testing"
)

func TestConflicts(t *testing.T) {
	t.Run("sub1", func(t *testing.T) {
		t.Parallel()
		RunTestConflicts(t, "main1", 5)
	})
	t.Run("sub2", func(t *testing.T) {
		t.Parallel()
		RunTestConflicts(t, "main2", 5)
	})
	t.Run("sub3", func(t *testing.T) {
		t.Parallel()
		RunTestConflicts(t, "main3", 5)
	})
}
