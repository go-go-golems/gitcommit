package ticket

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var TicketCmd *cobra.Command

var initOnce sync.Once
var initErr error

func Init() error {
	initOnce.Do(func() {
		if err := InitShowCmd(); err != nil {
			initErr = err
			return
		}
		TicketCmd = TicketShowCmd
	})
	return initErr
}

func must(err error) {
	if err != nil {
		panic(errors.Wrap(err, "init"))
	}
}
