package task

import (
	"time"
)

type BoolLocker bool

func (m *BoolLocker) LockTask(taskName string, execTimeout time.Duration) (success bool) {
	if *m == true {
		return
	}

	success = true
	*m = true
	return
}

func (m *BoolLocker) UnlockTask(taskName string) (success bool) {
	success = true
	*m = false
	return
}
