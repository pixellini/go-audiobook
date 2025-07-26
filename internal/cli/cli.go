package cli

import (
	"context"
	"fmt"

	"github.com/pixellini/go-audiobook/internal/app"
	"github.com/pixellini/go-audiobook/internal/flags"
)

func Run() error {
	return RunContext(context.Background())
}

func RunContext(ctx context.Context) error {
	f := flags.New()

	app, err := app.NewWithFlags(f)
	if err != nil {
		return fmt.Errorf("error happened on create", err)
	}

	err = app.Run()
	if err != nil {
		return fmt.Errorf("error happened on run", err)
	}

	return nil
}
