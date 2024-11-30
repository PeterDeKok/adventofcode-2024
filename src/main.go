package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/env"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/exit"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

func main() {
	defer exit.PanicToExit()

	styles.SetTerminalBackground()
	defer styles.ResetTerminalBackground()

	e := env.MustGet()

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	l := logger.MustInit()

	// TODO Should add ctx as first argument
	remoteClient := remote.CreateClient(&remote.ClientConfig{
		Base:                 "https://adventofcode.com/",
		PuzzlesDir:           e.PuzzlesDir,
		SessionCookieValue:   e.SessionCookieValue,
		SessionCookieExpires: e.SessionCookieExpires,
	})

	// TODO Should move ctx to first argument instead of config
	// TODO Move most values to env and retrieve in create function
	manager := manage.Create(&manage.AdventOfCodeConfig{
		Ctx:        ctx,
		Logger:     l,
		Remote:     remoteClient,
		PuzzlesDir: e.PuzzlesDir,
		FirstDay:   "2024-12-01",
		NrDays:     25,
		TZ:         "Europe/Amsterdam",
	})

	manager.MustLoadLocal()

	program := tui2.Start(ctx, manager, remoteClient)

	if _, err := program.Run(); err != nil {
		log.Error("runnign program failed", "err", err)
		styles.ResetTerminalBackground()
		os.Exit(1)
	}

	fmt.Println("Bye!")
}
