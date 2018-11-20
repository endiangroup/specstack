package metadata

import (
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func mockUids(t *testing.T, count int) (output []uuid.UUID) {
	for i := 0; i < count; i++ {
		uid := uuid.NewV5(uuid.NamespaceOID, fmt.Sprintf("%d", i))
		output = append(output, uid)
	}
	return
}
