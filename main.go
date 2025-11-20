package main

import (
	//"fmt"
	"log"
	"os"

	"github.com/atmetz/rssgator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	//fmt.Println("Working on commands")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	currentState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		command: make(map[string]func(*state, command) error),
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
	}

	cmds.register("login", handlerLogins)
	err = cmds.run(currentState, command{Name: os.Args[1], Args: os.Args[2:]})

	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}
