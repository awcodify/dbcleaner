package dbcleaner

import (
	"errors"
	"testing"

	"github.com/khaiql/dbcleaner/engine"
	"github.com/stretchr/testify/mock"
)

func TestClean(t *testing.T) {
	cleaner := New()
	mockEngine := &engine.MockEngine{}
	mockEngine.On("Truncate", mock.AnythingOfType("string")).Return(nil)
	mockEngine.On("Close").Return(nil)

	cleaner.SetEngine(mockEngine)

	t.Run("TestNothingLock", func(t *testing.T) {
		cleaner.Clean("table1", "table2")
		mockEngine.AssertNumberOfCalls(t, "Truncate", 2)
		mockEngine.AssertCalled(t, "Truncate", "table1")
		mockEngine.AssertCalled(t, "Truncate", "table2")
	})

	t.Run("TestLockAndThenUnlock", func(t *testing.T) {
		tbName := "lock_table"
		cleaner.Acquire(tbName)
		err := cleaner.Clean(tbName)
		if err != nil {
			t.Error(err.Error())
		}

		mockEngine.AssertCalled(t, "Truncate", tbName)
	})

	t.Run("TestClose", func(t *testing.T) {
		cleaner.Close()
		mockEngine.AssertCalled(t, "Close")
	})

	t.Run("TestTruncateError", func(t *testing.T) {
		errorTruncateMock := &engine.MockEngine{}
		errorTruncateMock.On("Truncate", mock.AnythingOfType("string")).Return(errors.New("Truncate error"))

		cleaner.SetEngine(errorTruncateMock)
		err := cleaner.Clean("error_table")

		if err.Error() != "Truncate error" {
			t.Error("Error mismatch")
		}
	})
}
