package cmd

// version is set at build time via:
//   go build -ldflags "-X github.com/lfry610/awsaudit/cmd.version=v0.1.0"
// Left as "dev" for local builds. This is the hook the "Shipping it"
// stage will use for versioned releases later — nothing to do with it
// yet beyond printing it.
var version = "dev"

// versionCmd is your first command to write end-to-end. It's
// deliberately simple so you can focus on the cobra.Command shape
// itself rather than any business logic.
//
// TODO(you):
//  0. Add "github.com/spf13/cobra" to the import above.
//  1. Build a *cobra.Command (Use: "version", Short: "...", RunE or Run)
//     that prints the version string above.
//  2. Register it with the root command: rootCmd.AddCommand(versionCmd)
//     — do this in an init() here, same pattern as cobra expects every
//     subcommand file to follow (keeps root.go from needing to know
//     about every subcommand that exists).
//
// Once this builds, `go run . version` should print "dev".
