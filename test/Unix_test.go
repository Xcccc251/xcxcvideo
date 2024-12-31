package test

import (
	"fmt"
	"testing"
	"time"
)

func TestId(t *testing.T) {
	unix := time.Now().Unix()
	fmt.Println(unix)
}
