package root

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"

	latest "github.com/docker/cagent/pkg/config/v1"
	"github.com/docker/cagent/pkg/server"
	"github.com/docker/cagent/pkg/session"
	"github.com/docker/cagent/pkg/teamloader"
)

var (
	listenAddr string
	sessionDb  string
	runConfig  latest.RuntimeConfig
)

// NewWebCmd creates a new web command
func NewApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api <agent-file>|<agents-dir>",
		Short: "Start the API server",
		Long:  `Start the API server that exposes the agent via an HTTP API`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHttp(cmd, false, args)
		},
	}

	cmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", ":8080", "Address to listen on")
	cmd.PersistentFlags().StringVarP(&sessionDb, "session-db", "s", "session.db", "Path to the session database")
	cmd.PersistentFlags().StringSliceVar(&runConfig.EnvFiles, "env-from-file", nil, "Set environment variables from file")
	addGatewayFlags(cmd)

	return cmd
}

func runHttp(cmd *cobra.Command, autoRunTools bool, args []string) error {
	ctx := cmd.Context()
	agentsPath := args[0]

	logger := newLogger()

	ln, err := server.Listen(ctx, listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenAddr, err)
	}
	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	if _, ok := ln.(*net.TCPListener); ok {
		logger.Info("Listening on http://localhost" + listenAddr)
	} else {
		logger.Info("Listening on " + listenAddr)
	}

	logger.Debug("Starting server", "agents", agentsPath, "debug_mode", debugMode)

	sessionStore, err := session.NewSQLiteSessionStore(sessionDb)
	if err != nil {
		return fmt.Errorf("failed to create session store: %w", err)
	}

	var opts []server.Opt

	stat, err := os.Stat(agentsPath)
	if err != nil {
		return fmt.Errorf("failed to stat agents path: %w", err)
	}
	if stat.IsDir() {
		opts = append(opts, server.WithAgentsDir(agentsPath))
	}

	teams, err := teamloader.LoadTeams(ctx, agentsPath, runConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to load teams: %w", err)
	}
	defer func() {
		for _, team := range teams {
			if err := team.StopToolSets(); err != nil {
				logger.Error("Failed to stop tool sets", "error", err)
			}
		}
	}()

	if autoRunTools {
		opts = append(opts, server.WithAutoRunTools(true))
	}

	s := server.New(logger, sessionStore, runConfig, teams, opts...)
	return s.Serve(ctx, ln)
}
