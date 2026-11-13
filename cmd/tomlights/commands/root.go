package commands

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/altessa-s/go-atlas/core/runtime/appinfo"

	"github.com/kitdoo/tomlights/cmd/tomlights/commands/server"
	"github.com/kitdoo/tomlights/cmd/tomlights/commands/version"
)

// command wraps [cobra.Command] to add tomlights-specific configuration
// such as subcommand registration and output stream setup.
type command struct {
	*cobra.Command
}

// New builds the root [cobra.Command] with all subcommands (server, version)
// already registered. The returned command is ready to execute.
func New() *cobra.Command {
	root := &command{
		Command: &cobra.Command{
			Use:   "tomlights",
			Short: "tomlights service",
			Long: heredoc.Doc(`
				tomlights is a Uber Fx / Cobra service skeleton built on go-atlas.

				It boots MongoDB, TLS, metrics, tracing, and health-check
				infrastructure, then exposes it through gRPC and HTTP servers with the
				standard probe/metrics endpoints. Domain logic is added on top.
			`),
			Example: heredoc.Doc(`
				# Start the server
				tomlights server run

				# Start with a specific configuration file
				tomlights server run --config /path/to/config.yaml

				# Show version information
				tomlights version

				# Show detailed help for server commands
				tomlights server --help
			`),
			SilenceUsage: true,
			Version:      appinfo.Version,
		},
	}

	root.configure()

	return root.Command
}

// configure wires stderr/stdout, enables typo suggestions, sets the
// version template, and registers the version and server subcommands.
func (c *command) configure() {
	c.SetErr(os.Stderr)
	c.SetOut(os.Stdout)
	c.SuggestionsMinimumDistance = 1
	c.SetVersionTemplate(heredoc.Doc(`
		{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
	`))
	c.AddCommand(version.New())
	c.AddCommand(server.New())
}

// Run creates a fresh root command, sets the given CLI arguments, and
// executes it. It is the main entry point called by the binary's main
// function, and can also be called directly in integration tests.
func Run(args []string) error {
	cmd := New()
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}
