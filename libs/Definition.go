package libs

import (
	"fmt"
	"time"
)

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type Limiter interface {
	Take() bool
	Put()
}

type LimitAndWaitFunc func() (int32, int32)

type CallError struct {
	Code    int
	Msg     string
	IsRetry bool
}

func (ce *CallError) Error() string {
	if ce == nil {
		return ""
	}
	return fmt.Sprintf("code:%d msg:%s", ce.Code, ce.Msg)
}
