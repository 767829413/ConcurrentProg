package tmp

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestOp(t *testing.T) {
	op()
}

func throwErr() error {
	return &MyErr{}
}

func op() {
	err := throwErr()
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case *MyErr:
			fmt.Println("呵呵呵", err)
		default:
			// unknown error
			fmt.Println("哈哈哈")
		}
	}
}

type MyErr struct {
	error
}

func (m *MyErr) Error() string {
	return "自己的错误"
}
