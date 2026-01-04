package preflight

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/gitcommit/pkg/app"
	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	gitcommit_layers "github.com/go-go-golems/gitcommit/pkg/layers"
	"github.com/go-go-golems/gitcommit/pkg/validate"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	glazed_layers "github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var CheckCmd *cobra.Command

type PreflightSettings struct {
	TicketOverride string `glazed.parameter:"ticket"`
	RequireStaged  bool   `glazed.parameter:"require-staged"`
	UseDocmgr      bool   `glazed.parameter:"docmgr"`
	AllowNoise     bool   `glazed.parameter:"allow-noise"`
}

type PreflightCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = &PreflightCommand{}

func NewPreflightCommand() (*PreflightCommand, error) {
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
	requireStaged := parameters.NewParameterDefinition(
		"require-staged",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Require staged files"),
		parameters.WithDefault(true),
	)
	useDocmgr := parameters.NewParameterDefinition(
		"docmgr",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Require docmgr to be initialized and ticket to exist"),
		parameters.WithDefault(true),
	)
	allowNoise := parameters.NewParameterDefinition(
		"allow-noise",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Allow common noise files (dist/, node_modules/, .env, etc.) to be staged"),
		parameters.WithDefault(false),
	)

	cmdDesc := cmds.NewCommandDescription(
		"preflight",
		cmds.WithShort("Validate repo state before committing"),
		cmds.WithLong("Runs the same safety checks as `gitcommit commit`: ticket detection, staged files, noise policy, and optional docmgr readiness."),
		cmds.WithFlags(ticketFlag, requireStaged, useDocmgr, allowNoise),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &PreflightCommand{CommandDescription: cmdDesc}, nil
}

func (c *PreflightCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *glazed_layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &PreflightSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode preflight settings")
	}

	repoRoot, err := git.RepoRoot(ctx, repoSettings.RepoPath)
	if err != nil {
		return err
	}

	ticketRes, err := app.ResolveTicket(ctx, repoRoot, settings.TicketOverride)
	if err != nil {
		return err
	}

	stagedFiles, err := git.StagedFiles(ctx, repoRoot)
	if err != nil {
		return err
	}
	if settings.RequireStaged && len(stagedFiles) == 0 {
		return errors.New("no staged files; stage changes first (git add ...)")
	}

	noise := []validate.NoiseFinding(nil)
	if !settings.AllowNoise {
		noise = validate.FindNoise(stagedFiles)
		if len(noise) > 0 {
			var b strings.Builder
			b.WriteString("refusing to proceed due to common noise files (use --allow-noise to override):\n")
			for _, n := range noise {
				b.WriteString("- ")
				b.WriteString(n.Path)
				b.WriteString(" (")
				b.WriteString(n.Reason)
				b.WriteString(")\n")
			}
			return errors.New(b.String())
		}
	}

	docmgrOK := false
	if settings.UseDocmgr {
		if !docmgr.IsAvailable() {
			return errors.New("docmgr not found in PATH (use --docmgr=false to disable)")
		}
		// docmgr config lives at repo root.
		if !hasDocmgrConfig(repoRoot) {
			return errors.New("docmgr not initialized (missing .ttmp.yaml); run: gitcommit docmgr init")
		}
		exists, err := docmgr.TicketExists(ctx, repoRoot, ticketRes.TicketID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("docmgr ticket %s not found; create it with: gitcommit docmgr ticket create --ticket %s --title \"...\" --topics ...", ticketRes.TicketID, ticketRes.TicketID)
		}
		docmgrOK = true
	}

	row := types.NewRow(
		types.MRP("repo_root", repoRoot),
		types.MRP("ticket", ticketRes.TicketID),
		types.MRP("ticket_source", ticketRes.Source),
		types.MRP("staged_files", len(stagedFiles)),
		types.MRP("docmgr", settings.UseDocmgr),
		types.MRP("docmgr_ok", docmgrOK),
		types.MRP("allow_noise", settings.AllowNoise),
		types.MRP("noise_findings", len(noise)),
	)
	return gp.AddRow(ctx, row)
}

func hasDocmgrConfig(repoRoot string) bool {
	_, err := os.Stat(filepath.Join(repoRoot, ".ttmp.yaml"))
	return err == nil
}

func InitCheckCmd() error {
	glazedCmd, err := NewPreflightCommand()
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

	CheckCmd = cobraCmd
	return nil
}
