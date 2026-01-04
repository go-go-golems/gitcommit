package commit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/gitcommit/pkg/app"
	"github.com/go-go-golems/gitcommit/pkg/commitmsg"
	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	gitcommit_layers "github.com/go-go-golems/gitcommit/pkg/layers"
	"github.com/go-go-golems/gitcommit/pkg/validate"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	glazed_layers "github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var CommitOnceCmd *cobra.Command

type CommitSettings struct {
	Message        string   `glazed.parameter:"message"`
	Body           []string `glazed.parameter:"body"`
	TicketOverride string   `glazed.parameter:"ticket"`
	UseDocmgr      bool     `glazed.parameter:"docmgr"`
	DryRun         bool     `glazed.parameter:"dry-run"`
	AllowNoise     bool     `glazed.parameter:"allow-noise"`
}

type CommitCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &CommitCommand{}

func NewCommitCommand() (*CommitCommand, error) {
	repoLayer, err := gitcommit_layers.NewRepositoryLayer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create repository layer")
	}
	repoLayerExisting, err := gitcommit_layers.WrapAsExistingCobraFlagsLayer(repoLayer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap repository layer as existing flags layer")
	}

	msgFlag := parameters.NewParameterDefinition(
		"message",
		parameters.ParameterTypeString,
		parameters.WithHelp("Commit summary (ticket prefix will be ensured)"),
		parameters.WithDefault(""),
		parameters.WithShortFlag("m"),
	)
	bodyFlag := parameters.NewParameterDefinition(
		"body",
		parameters.ParameterTypeStringList,
		parameters.WithHelp("Commit body paragraph (repeatable)"),
		parameters.WithDefault([]string{}),
	)
	ticketFlag := parameters.NewParameterDefinition(
		"ticket",
		parameters.ParameterTypeString,
		parameters.WithHelp("Ticket ID override (defaults to env/branch detection)"),
		parameters.WithDefault(""),
	)
	docmgrFlag := parameters.NewParameterDefinition(
		"docmgr",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Update docmgr changelog for the ticket"),
		parameters.WithDefault(true),
	)
	dryRunFlag := parameters.NewParameterDefinition(
		"dry-run",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Print what would happen without committing"),
		parameters.WithDefault(false),
	)
	allowNoiseFlag := parameters.NewParameterDefinition(
		"allow-noise",
		parameters.ParameterTypeBool,
		parameters.WithHelp("Allow common noise files (dist/, node_modules/, .env, etc.) to be committed"),
		parameters.WithDefault(false),
	)

	cmdDesc := cmds.NewCommandDescription(
		"commit",
		cmds.WithShort("Commit with ticket prefix and optional docmgr changelog update"),
		cmds.WithLong("Creates a git commit using the staged files, ensuring the ticket prefix in the summary. Optionally updates the ticket changelog via docmgr."),
		cmds.WithFlags(msgFlag, bodyFlag, ticketFlag, docmgrFlag, dryRunFlag, allowNoiseFlag),
		cmds.WithLayersList(repoLayerExisting),
	)

	return &CommitCommand{CommandDescription: cmdDesc}, nil
}

func (c *CommitCommand) Run(ctx context.Context, parsedLayers *glazed_layers.ParsedLayers) error {
	repoSettings, err := gitcommit_layers.GetRepositorySettings(parsedLayers)
	if err != nil {
		return err
	}

	settings := &CommitSettings{}
	if err := parsedLayers.InitializeStruct(glazed_layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to decode commit settings")
	}

	if strings.TrimSpace(settings.Message) == "" {
		return errors.New("missing --message")
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
	if len(stagedFiles) == 0 {
		return errors.New("no staged files; stage changes first (git add ...)")
	}

	if !settings.AllowNoise {
		noise := validate.FindNoise(stagedFiles)
		if len(noise) > 0 {
			var b strings.Builder
			b.WriteString("refusing to commit common noise files (use --allow-noise to override):\n")
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

	finalMessage := commitmsg.EnsureTicketPrefix(ticketRes.TicketID, settings.Message)

	if settings.DryRun {
		fmt.Fprintf(os.Stdout, "repo:   %s\n", repoRoot)
		fmt.Fprintf(os.Stdout, "ticket: %s (%s)\n", ticketRes.TicketID, ticketRes.Source)
		fmt.Fprintf(os.Stdout, "msg:    %s\n", finalMessage)
		fmt.Fprintf(os.Stdout, "files:  %d staged\n", len(stagedFiles))
		fmt.Fprintf(os.Stdout, "docmgr: %v\n", settings.UseDocmgr)
		return nil
	}

	if settings.UseDocmgr {
		if !docmgr.IsAvailable() {
			return errors.New("docmgr not found in PATH (use --docmgr=false to disable)")
		}
		if _, err := os.Stat(filepath.Join(repoRoot, ".ttmp.yaml")); err != nil {
			return errors.Wrap(err, "docmgr not initialized (missing .ttmp.yaml); run: gitcommit docmgr init")
		}
		exists, err := docmgr.TicketExists(ctx, repoRoot, ticketRes.TicketID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Errorf("docmgr ticket %s not found; create it with: gitcommit docmgr ticket create --ticket %s --title \"...\" --topics ...", ticketRes.TicketID, ticketRes.TicketID)
		}
	}

	if err := git.Commit(ctx, repoRoot, finalMessage, settings.Body); err != nil {
		return err
	}

	hash, err := git.HeadHash(ctx, repoRoot)
	if err != nil {
		return err
	}

	if settings.UseDocmgr {
		fileNotes := make([]docmgr.FileNote, 0, len(stagedFiles))
		for _, f := range stagedFiles {
			fileNotes = append(fileNotes, docmgr.FileNote{
				Path: filepath.Join(repoRoot, f),
				Note: "Changed in commit " + hash,
			})
		}

		entry := fmt.Sprintf("Commit %s: %s", hash, finalMessage)
		if err := docmgr.ChangelogUpdate(ctx, repoRoot, ticketRes.TicketID, entry, fileNotes); err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "docmgr updated; commit docs changes separately (e.g. git add ttmp/... && git commit -m \"Docs: update changelog\")")
	}

	fmt.Fprintln(os.Stdout, hash)
	return nil
}

func InitCommitCmd() error {
	glazedCmd, err := NewCommitCommand()
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

	CommitOnceCmd = cobraCmd
	return nil
}
