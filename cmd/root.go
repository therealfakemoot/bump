/*
Copyright © 2019 Guilhem Lettron <guilhem@barpilot.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/therealfakemoot/bump/pkg/git"
	"github.com/therealfakemoot/bump/pkg/semver"
)

type LogContextKey string

const LogKey LogContextKey = "logger"

var (
	ErrDirtyRepo         = errors.New("repo is dirty")
	ErrCommandNotMatched = errors.New("command not matched")
)

var (
	allowDirty bool
	currentTag string
	latestTag  bool
	dryRun     bool
)

var rootCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump version",
	Long:  ``,

	SilenceUsage: true,

	PersistentPreRunE: preRun,
}

func Execute() {
	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConfig.DisableStacktrace = true
	log := zap.Must(logConfig.Build())
	defer log.Sync()

	g, err := git.New()
	if err != nil {
		log.Fatal("error loading git repo", zap.Error(err))
	}

	ctx := context.WithValue(
		context.Background(),
		LogKey,
		log,
	)
	ctx = context.WithValue(ctx, "git", g)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error("error executing rootCmd", zap.Error(err))
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVar(&allowDirty, "allow-dirty", false, "allow usage of bump on dirty git")
	rootCmd.PersistentFlags().BoolVar(&latestTag, "latest-tag", true, "use latest tag, prompt tags if false")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Don't touch git repository")
}

func preRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	log := ctx.Value(LogKey).(*zap.Logger)

	g := ctx.Value("git").(*git.Git)

	if !allowDirty {
		if g.IsDirty() {
			log.Error("repo is dirty", zap.Error(ErrDirtyRepo))
		}
	}

	tags, err := g.Tags()
	if err != nil {
		log.Error("unable to fetch tags", zap.Error(err))

		return err
	}

	log = log.With(zap.Strings("tags", tags))

	if !latestTag {
		prompt := promptui.Select{
			Label: "Select Previous tag",
			Items: tags,
		}

		_, currentTag, err = prompt.Run()

		if err != nil {
			log.Error("prompt threw an error", zap.Error(err))
			return err
		}
	} else {
		currentTag, err = semver.Latest(tags)
		if err != nil {
			log.Error("semver could not fetch latest tag", zap.Error(err))
			return err
		}
	}

	log.Info("tag chosen", zap.String("current tag", currentTag))

	return nil
}
