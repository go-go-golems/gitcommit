package ticket

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

var ExistsCmd *cobra.Command

type ExistsSettings struct {
	TicketOverride string `glazed.parameter:"ticket"`
}

type ExistsCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &ExistsCommand{}

func NewExistsCommand() (*ExistsCommand, error) {
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

	cmdDesc := cmds.NewCommandDescription(
		"exists",
		cmds.WithShort("Exit 0 if the docmgr ticket exists"),
		cmds.WithLong("Checks whether a docmgr ticket exists; prints the resolved ticket ID if it does."),
		cmds.WithFlags(ticketFlag),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &ExistsCommand{CommandDescription: cmdDesc}, nil
}

func (c *ExistsCommand) Run(ctx context.Context, parsedLayers *glazed_layers.ParsedLayers) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &ExistsSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode exists settings")
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

	exists, err := docmgr.TicketExists(ctx, repoRoot, ticketRes.TicketID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("docmgr ticket %s not found", ticketRes.TicketID)
	}

	_, _ = fmt.Fprintln(os.Stdout, ticketRes.TicketID)
	return nil
}

func InitExistsCmd() error {
	glazedCmd, err := NewExistsCommand()
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

	ExistsCmd = cobraCmd
	return nil
}
