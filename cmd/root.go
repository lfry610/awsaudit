package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyProfile = "profile"
	keyRegion  = "region"
	keyOutput  = "output"
)

// cfgFile holds the --config flag value (path to an optional config file).
// It's a package-level var because cobra's flag parsing and viper's config
// loading happen in two separate callbacks (init() and initConfig()) that
// both need it.
var cfgFile string

// rootCmd is the base command. Every other command (version, scan, ...)
// gets attached to this one via rootCmd.AddCommand in their own init().
var rootCmd = &cobra.Command{
	Use:   "awsaudit",
	Short: "awsaudit scans an AWS account for hygiene problems",
	Long: `awsaudit is a CLI that points at an AWS account and reports
	common hygiene issues: missing tags, overly open security groups,
	public/unencrypted storage, stale IAM keys, and idle billable resources.`,
	// SilenceUsage/SilenceErrors: without these, cobra prints the full
	// usage block on every error returned from a RunE, which is noisy
	// for things like "no AWS credentials found". We want our own
	// error formatting (see Execute below) to be what the user sees.
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute is the single entry point called from main.go. Keeping all
// error handling and os.Exit here (rather than scattering os.Exit calls
// through subcommands) means exit codes stay predictable: this is the
// "exit codes other tools can rely on" piece from the project goals.
// 0 = success, 1 = error, 2 = findings were reported (we'll wire that up
// once there's an actual scanner — for now just 0/1).
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file (default $HOME/awsaudit.yaml)")
	pf.String(keyProfile, "", "AWS named profile to use")
	pf.String(keyRegion, "", "AWS region to target")
	pf.StringP(keyOutput, "o", "table", "output format: table|json")

	// --- Global/persistent flags ---
	// PersistentFlags() (vs. Flags()) means these are inherited by every
	// subcommand (version, and later scan/scan-tags/scan-sg/...), not
	// just rootCmd itself.
	//
	// TODO(you): declare the global flags this tool needs. At minimum:
	//   --config string   path to a config file (bind to &cfgFile)
	//   --profile string  AWS named profile to use
	//   --region string   AWS region to target
	//   --output string   output format, e.g. "table" or "json" (short: -o)
	//
	// Use rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "...")
	// for --config, and StringVarP for ones that want a short flag like -o.
	//
	// Why flags for things like --profile/--region instead of just reading
	// env vars directly? Because viper (wired in initConfig) will let a
	// user set these via flag, env var, or config file, with flags always
	// winning. That three-way precedence is the thing this stage is
	// teaching — see the TODO in initConfig below.
}

// initConfig wires viper up to read configuration from (in order of
// precedence, highest first): command-line flags > environment
// variables > a config file > defaults.
//
// viper isn't a dependency yet since nothing imports it. Run
// `go get github.com/spf13/viper` and add it to the import block
// above when you start filling this in.
//
// TODO(you): implement this. Rough shape:
//
//  1. If cfgFile is set (user passed --config), tell viper to use that
//     exact file: viper.SetConfigFile(cfgFile).
//     Otherwise, give viper a name/type/search path to look for one
//     automatically, e.g. viper.SetConfigName("awsaudit"),
//     viper.SetConfigType("yaml"), viper.AddConfigPath("$HOME").
//
//  2. Enable env var support: viper.SetEnvPrefix("AWSAUDIT") and
//     viper.AutomaticEnv(), so e.g. AWSAUDIT_REGION=us-east-1 works.
//
//  3. Read the config file with viper.ReadInConfig(). A missing file
//     should NOT be a fatal error (config files are optional) — check
//     for viper.ConfigFileNotFoundError specifically and ignore just
//     that case; report any other error.
//
//  4. Bind each persistent flag you declared above to viper, e.g.
//     viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region")).
//     This is what makes the precedence chain work: later, instead of
//     reading the --region flag directly, subcommands should call
//     viper.GetString("region") and get the flag value, or the env var,
//     or the config file value, or a default — in that priority order,
//     without the subcommand needing to care which one fired.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(home)
		}
		viper.SetConfigName("awsaudit")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("AWSAUDIT")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			fmt.Fprintf(os.Stderr, "Error: reading config: %v\n", err)
			os.Exit(1)
		}
	}

	viper.BindPFlag(keyProfile, rootCmd.PersistentFlags().Lookup(keyProfile))
	viper.BindPFlag(keyRegion, rootCmd.PersistentFlags().Lookup(keyRegion))
	viper.BindPFlag(keyOutput, rootCmd.PersistentFlags().Lookup(keyOutput))
}
