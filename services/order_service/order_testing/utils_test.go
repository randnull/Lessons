package order_testing

import (
	"github.com/randnull/Lessons/pkg/logger"
	"testing"
)

func TestInitLogger(t *testing.T) {
	err := logger.InitLogger()

	if err != nil {
		t.Errorf("%v", err)
	}
}
