package ticket

import (
	"context"

	"github.com/go-go-golems/gitcommit/pkg/app"
	"github.com/go-go-golems/gitcommit/pkg/git"
	gitcommit_layers "github.com/go-go-golems/gitcommit/pkg/layers"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	glazed_layers "github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var TicketShowCmd *cobra.Command

type TicketShowSettings struct {
	Ticket string `glazed.parameter:"ticket"`
}

type TicketShowCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = &TicketShowCommand{}

func NewTicketShowCommand() (*TicketShowCommand, error) {
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
		"ticket",
		cmds.WithShort("Show detected ticket ID"),
		cmds.WithLong("Detect the ticket ID from --ticket, GITCOMMIT_TICKET, or the current git branch name."),
		cmds.WithFlags(ticketFlag),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &TicketShowCommand{CommandDescription: cmdDesc}, nil
}

func (c *TicketShowCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *glazed_layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &TicketShowSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode ticket settings")
	}

	repoRoot, err := git.RepoRoot(ctx, repoSettings.RepoPath)
	if err != nil {
		return err
	}

	res, err := app.ResolveTicket(ctx, repoRoot, settings.Ticket)
	if err != nil {
		return err
	}

	row := types.NewRow(
		types.MRP("ticket", res.TicketID),
		types.MRP("source", res.Source),
		types.MRP("branch", res.Branch),
		types.MRP("repo_root", repoRoot),
	)
	return gp.AddRow(ctx, row)
}

func InitShowCmd() error {
	glazedCmd, err := NewTicketShowCommand()
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

	TicketShowCmd = cobraCmd
	return nil
}
