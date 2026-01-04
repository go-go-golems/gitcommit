package docmgr

import (
	"context"

	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	gitcommit_layers "github.com/go-go-golems/gitcommit/pkg/layers"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	glazed_layers "github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var InitCmd *cobra.Command

type InitCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &InitCommand{}

func NewInitCommand() (*InitCommand, error) {
	repoLayer, err := gitcommit_layers.NewRepositoryLayer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create repository layer")
	}
	repoLayerExisting, err := gitcommit_layers.WrapAsExistingCobraFlagsLayer(repoLayer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap repository layer as existing flags layer")
	}

	cmdDesc := cmds.NewCommandDescription(
		"init",
		cmds.WithShort("Initialize docmgr in the repository"),
		cmds.WithLong("Runs docmgr init --seed-vocabulary in the repository root."),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &InitCommand{CommandDescription: cmdDesc}, nil
}

func (c *InitCommand) Run(ctx context.Context, parsedLayers *glazed_layers.ParsedLayers) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	repoRoot, err := git.RepoRoot(ctx, repoSettings.RepoPath)
	if err != nil {
		return err
	}

	if !docmgr.IsAvailable() {
		return errors.New("docmgr not found in PATH")
	}

	return docmgr.Init(ctx, repoRoot)
}

func InitInitCmd() error {
	glazedCmd, err := NewInitCommand()
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

	InitCmd = cobraCmd
	return nil
}
