package docmgr

import (
	"context"
	"fmt"
	"os"

	"github.com/go-go-golems/gitcommit/pkg/app"
	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	gitcommit_layers "github.com/go-go-golems/gitcommit/pkg/layers"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	glazed_layers "github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var DoctorCmd *cobra.Command

type DoctorSettings struct {
	TicketOverride string `glazed.parameter:"ticket"`
	StaleAfter     int    `glazed.parameter:"stale-after"`
}

type DoctorCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &DoctorCommand{}

func NewDoctorCommand() (*DoctorCommand, error) {
	repoLayer, err := gitcommit_layers.NewRepositoryLayer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create repository layer")
	}
	repoLayerExisting, err := gitcommit_layers.WrapAsExistingCobraFlagsLayer(repoLayer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap repository layer as existing flags layer")
	}

	ticketFlag := parameters.NewParameterDefinition(
		"ticket",
		parameters.ParameterTypeString,
		parameters.WithHelp("Ticket ID override (defaults to env/branch detection)"),
		parameters.WithDefault(""),
	)
	staleAfter := parameters.NewParameterDefinition(
		"stale-after",
		parameters.ParameterTypeInteger,
		parameters.WithHelp("Stale-after threshold in days"),
		parameters.WithDefault(30),
	)

	cmdDesc := cmds.NewCommandDescription(
		"doctor",
		cmds.WithShort("Run docmgr doctor for the ticket"),
		cmds.WithLong("Runs docmgr doctor --ticket <ticket> and prints the report."),
		cmds.WithFlags(ticketFlag, staleAfter),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &DoctorCommand{CommandDescription: cmdDesc}, nil
}

func (c *DoctorCommand) Run(ctx context.Context, parsedLayers *glazed_layers.ParsedLayers) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &DoctorSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode doctor settings")
	}

	repoRoot, err := git.RepoRoot(ctx, repoSettings.RepoPath)
	if err != nil {
		return err
	}

	if !docmgr.IsAvailable() {
		return errors.New("docmgr not found in PATH")
	}

	ticketRes, err := app.ResolveTicket(ctx, repoRoot, settings.TicketOverride)
	if err != nil {
		return err
	}

	out, err := docmgr.Doctor(ctx, repoRoot, ticketRes.TicketID, settings.StaleAfter)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(os.Stdout, out)
	return nil
}

func InitDoctorCmd() error {
	glazedCmd, err := NewDoctorCommand()
	if err != nil {
		return err
	}

	cobraCmd, err := cli.BuildCobraCommand(
		glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			MiddlewaresFunc: cli.CobraCommandDefaultMiddlewares,
		}),
	)
	if err != nil {
		return err
	}

	DoctorCmd = cobraCmd
	return nil
}
