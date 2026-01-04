package docmgr

import (
	"sync"

	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds/docmgr/ticket"
	"github.com/spf13/cobra"
)

var DocmgrCmd = &cobra.Command{
	Use:   "docmgr",
	Short: "Helpers around docmgr (init, ticket create, doctor)",
}

var initOnce sync.Once
var initErr error

func Init() error {
	initOnce.Do(func() {
		if err := InitInitCmd(); err != nil {
			initErr = err
			return
		}
		if err := InitDoctorCmd(); err != nil {
			initErr = err
			return
		}
		if err := ticket.Init(); err != nil {
			initErr = err
			return
		}

		DocmgrCmd.AddCommand(
			InitCmd,
			DoctorCmd,
			ticket.TicketCmd,
		)
	})
	return initErr
}
