package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/ivanklee86/argonap/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Build information (injected by goreleaser).
	version = "dev"
)

const (
	defaultConfigFilename = "argonap"
	envPrefix             = "ARGONAP"
)

// main function.
func main() {
	command := NewRootCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	argonap := cli.New()

	cmd := &cobra.Command{
		Use:     "argonap",
		Short:   "Give ArgoCD a lil' nap.",
		Long:    "A CLI to provision sync windows at scale.",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			argonap.Out = cmd.OutOrStdout()
			argonap.Err = cmd.ErrOrStderr()

			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(argonap.Out, cmd.UsageString())
		},
	}

	cmd.PersistentFlags().StringVar(&argonap.Config.ServerAddr, "server-address", "", "ArgoCD server address")
	cmd.PersistentFlags().BoolVar(&argonap.Config.Insecure, "insecure", false, "Don't validate SSL certificate on client request")
	cmd.PersistentFlags().StringVar(&argonap.Config.AuthToken, "auth-token", "", "JWT Authentication Token")
	cmd.PersistentFlags().StringVar(&argonap.Config.ProjectName, "name", "", "Project name to update.  If specified, label filtering will not apply.")
	cmd.PersistentFlags().StringSliceVar(&argonap.Config.LabelsAsStrings, "label", []string{}, "Labels to filter projects on in format 'key=value'.  Can be used multiple times.")

	cmd.AddCommand(NewClearCommand(argonap))
	cmd.AddCommand(NewSetCommand(argonap))

	return cmd
}

func NewClearCommand(argonap *cli.argonap) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear SyncWindows.",
		Long:  "Clear SyncWindows on all AppProjects.",
		Run: func(cmd *cobra.Command, args []string) {
			argonap.Connect()
			argonap.ClearSyncWindows()
		},
	}

	return cmd
}

func NewSetCommand(argonap *cli.argonap) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set SyncWindows.",
		Long:  "Set SyncWindows from file",
		Run: func(cmd *cobra.Command, args []string) {
			argonap.Connect()
			argonap.SetSyncWindows()
		},
	}
	cmd.PersistentFlags().StringVar(&argonap.Config.SyncWindowsFile, "file", "", "Path to file with SyncWindows to configure")

	return cmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			if err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix)); err != nil {
				os.Exit(1)
			}
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				os.Exit(1)
			}
		}
	})
}
