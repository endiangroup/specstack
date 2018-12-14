package specification

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func Test_ASnapshotterCanSnapshot(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/a.feature": mockFeatureA,
			"features/b.feature": mockFeatureB,
			"features/i.feature": mockFeatureI,
		},
	)
	mockReadSourcer := &MockReadSourcer{}
	mockObjectHasher := &MockObjectHasher{}
	snapshotter := NewSnapshotter(mockReadSourcer, mockObjectHasher)

	ss, err := snapshotter.Snapshot(spec)
	require.Nil(t, err)
	snaptest.Snapshot(t, ss)
}
