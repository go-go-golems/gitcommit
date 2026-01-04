package commit

import (
	"sync"

	"github.com/spf13/cobra"
)

var CommitCmd *cobra.Command

var initOnce sync.Once
var initErr error

func Init() error {
	initOnce.Do(func() {
		if err := InitCommitCmd(); err != nil {
			initErr = err
			return
		}
		CommitCmd = CommitOnceCmd
	})
	return initErr
}
