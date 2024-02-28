package cmd

import (
	"go.uber.org/zap"

	"github.com/spf13/cobra"

	"github.com/therealfakemoot/bump/pkg/git"
	"github.com/therealfakemoot/bump/pkg/semver"
)

func inc(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	log := ctx.Value(LogKey).(*zap.Logger)

	log = log.With(
		zap.String("command path", cmd.CommandPath()),
		zap.String("current tag", currentTag),
	)

	version := semver.New(currentTag)

	switch cmd.CommandPath() {
	case "bump patch":
		version.IncPatch()
	case "bump minor":
		version.IncMinor()
	case "bump major":
		version.IncMajor()
	default:
		log.Error("execution failed", zap.Error(ErrCommandNotMatched))
		return ErrCommandNotMatched
	}

	log = log.With(zap.String("new tag", version.StringFull()))

	log.Debug("values")

	if !dryRun {
		g := ctx.Value("git").(*git.Git)

		if err := g.CreateTag(version.StringFull()); err != nil {
			log.Error("could not create tag", zap.Error(err))
			return err
		}
	}

	return nil
}
