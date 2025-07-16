package cli

import (
	"flag"
	"fmt"
	"time"

	"github.com/pixellini/go-audiobook/internal/app"
)

// CLIFlags represents command-line flags
type CLIFlags struct {
	ResetProgress   bool
	FinishAudiobook bool
}

// Execute is the main entry point for the CLI application
func Execute() error {
	start := time.Now()

	// Parse command-line flags
	flags := parseFlags()

	// Create and configure the application
	application, err := app.New()
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	// Run the application with the provided flags
	if err := application.Run(flags.ResetProgress, flags.FinishAudiobook); err != nil {
		return fmt.Errorf("application failed: %w", err)
	}

	fmt.Printf("Audiobook created! Total time: %v\n", time.Since(start))
	return nil
}

// parseFlags parses command-line flags and returns a CLIFlags struct
func parseFlags() CLIFlags {
	flags := CLIFlags{}

	flag.BoolVar(&flags.ResetProgress, "reset", false, "Reset the audiobook generation process")
	flag.BoolVar(&flags.FinishAudiobook, "finish", false, "Finish audiobook generation with currently processed chapters")
	flag.Parse()

	return flags
}
