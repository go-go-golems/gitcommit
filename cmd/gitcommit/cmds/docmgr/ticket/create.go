package ticket

import (
	"context"

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

var CreateCmd *cobra.Command

type CreateSettings struct {
	TicketOverride string   `glazed.parameter:"ticket"`
	Title          string   `glazed.parameter:"title"`
	Topics         []string `glazed.parameter:"topics"`
}

type CreateCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &CreateCommand{}

func NewCreateCommand() (*CreateCommand, error) {
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
	titleFlag := parameters.NewParameterDefinition(
		"title",
		parameters.ParameterTypeString,
		parameters.WithHelp("Ticket title"),
		parameters.WithDefault(""),
	)
	topicsFlag := parameters.NewParameterDefinition(
		"topics",
		parameters.ParameterTypeStringList,
		parameters.WithHelp("Topics for the ticket"),
		parameters.WithDefault([]string{"chat"}),
	)

	cmdDesc := cmds.NewCommandDescription(
		"create",
		cmds.WithShort("Create a docmgr ticket workspace"),
		cmds.WithLong("Creates a docmgr ticket workspace under ttmp/ and seeds it with index/tasks/changelog."),
		cmds.WithFlags(ticketFlag, titleFlag, topicsFlag),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &CreateCommand{CommandDescription: cmdDesc}, nil
}

func (c *CreateCommand) Run(ctx context.Context, parsedLayers *glazed_layers.ParsedLayers) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &CreateSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode create settings")
	}
	if settings.Title == "" {
		return errors.New("missing --title")
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

	return docmgr.CreateTicket(ctx, repoRoot, ticketRes.TicketID, settings.Title, settings.Topics)
}

func InitCreateCmd() error {
	glazedCmd, err := NewCreateCommand()
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

	CreateCmd = cobraCmd
	return nil
}
