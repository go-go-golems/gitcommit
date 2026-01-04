package preflight

import (
	"sync"

	"github.com/spf13/cobra"
)

var PreflightCmd *cobra.Command

var initOnce sync.Once
var initErr error

func Init() error {
	initOnce.Do(func() {
		if err := InitCheckCmd(); err != nil {
			initErr = err
			return
		}
		PreflightCmd = CheckCmd
	})
	return initErr
}
