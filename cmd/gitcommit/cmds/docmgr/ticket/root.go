package ticket

import (
	"sync"

	"github.com/spf13/cobra"
)

var TicketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Docmgr ticket helpers",
}

var initOnce sync.Once
var initErr error

func Init() error {
	initOnce.Do(func() {
		if err := InitCreateCmd(); err != nil {
			initErr = err
			return
		}
		if err := InitExistsCmd(); err != nil {
			initErr = err
			return
		}

		TicketCmd.AddCommand(
			CreateCmd,
			ExistsCmd,
		)
	})
	return initErr
}
