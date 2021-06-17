package base

import (
	"context"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer func() {
			t.Log("exit")
		}()
		for {
			select {
			case <-ctx.Done():
				t.Log("exit")
			default:
				time.Sleep(time.Second)
			}
		}
	}(ct, cancel)
	time.Sleep(10 * time.Second)
}
