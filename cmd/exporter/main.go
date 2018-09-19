package main

import (
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/LOG-ED/generic-exporter/pkg/generic"
)

func newExporterCommand() *cobra.Command {
	var address string
	var metricsPath string

	cmd := &cobra.Command{
		Use:   "exporter [OPTIONS]",
		Short: "A Prometheus exporter.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if address == "" {
				return fmt.Errorf("The listen address cannot be empty")
			}
			if metricsPath == "" {
				return fmt.Errorf("The metrics path cannot be empty")
			}
			
			return generic.Run(address, metricsPath)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&address, "listen-address", ":8080", "Address to listen on for HTTP requests")
	flags.StringVar(&metricsPath, "metrics-path", "/metrics", "Path to metrics endpoint")
	return cmd
}

func main() {
	log.Info("Starting Generic Exporter")
	cmd := newExporterCommand()
	if err := cmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}