package dispatcher

import (
	"fmt"
	"testing"

	"github.com/takawang/sugar"
)

// Test Cases =========================
func TestInit(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Test", func(logf sugar.Log) bool {
		loadRouteFile("testdata/route.json")
		return true
	})

	if s.IsFailed() {
		fmt.Println("the tests failed :/")
	}
}
