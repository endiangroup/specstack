package specification

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func newSnapshotOfMockSpec(t *testing.T, files map[string]string) Snapshot {
	spec := generateAndReadSpec(t, files)
	mockReadSourcer := &MockReadSourcer{}
	mockObjectHasher := &MockObjectHasherPlaintext{}
	snapshotter := NewSnapshotter(mockReadSourcer, mockObjectHasher)
	ss, err := snapshotter.Snapshot(spec)
	require.Nil(t, err)
	return ss
}

func Test_ASnapshotCanDiff(t *testing.T) {
	s0 := newSnapshotOfMockSpec(t,
		map[string]string{
			"features/a.feature": mockFeatureA,
			"features/b.feature": mockFeatureB,
		},
	)
	s1 := newSnapshotOfMockSpec(t,
		map[string]string{
			"features/b.feature": mockFeatureB,
			"features/c.feature": mockFeatureC,
		},
	)
	r, a := s0.Diff(s1)
	t.Run("Removed", func(t *testing.T) {
		snaptest.Snapshot(t, r)
	})
	t.Run("Added", func(t *testing.T) {
		snaptest.Snapshot(t, a)
	})
}
