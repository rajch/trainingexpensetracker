package memrepository

import (
	"testing"

	"github.com/rajch/trainingexpensetracker/internal/model"
)

func TestMemExpenseRepository(t *testing.T) {
	repo := NewMemExpenseRepository()
	model.RunRepositoryTests(t, repo)
}
