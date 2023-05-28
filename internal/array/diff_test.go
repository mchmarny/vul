package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayDiffs(t *testing.T) {
	newIDs := []int64{1, 2, 3, 4, 5}
	dbIDs := []int64{6, 7, 3, 4, 5}

	deletes := GetDiff(newIDs, dbIDs)
	assert.Contains(t, deletes, int64(6))
	assert.Contains(t, deletes, int64(7))

	inserts := GetDiff(dbIDs, newIDs)
	assert.Contains(t, inserts, int64(1))
	assert.Contains(t, inserts, int64(2))

	unique := Unique(newIDs, dbIDs)
	assert.Len(t, unique, 7)
	assert.Contains(t, unique, int64(1))
	assert.Contains(t, unique, int64(2))
	assert.Contains(t, unique, int64(3))
	assert.Contains(t, unique, int64(4))
	assert.Contains(t, unique, int64(5))
	assert.Contains(t, unique, int64(6))
	assert.Contains(t, unique, int64(7))
}
