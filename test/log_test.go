package test

import (
	"context"
	"testing"

	"go-spring.org/log"
)

func TestLog(t *testing.T) {
	log.Info(context.Background(), log.TagAppDef, "hello world")
}
