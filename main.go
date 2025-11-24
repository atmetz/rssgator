package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/atmetz/rssgator/internal/config"
	"github.com/atmetz/rssgator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Read config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	//	Connect to database

	dbURL := cfg.DBURL

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	// create state
	currentState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Create cmds struct to hold registered commands
	cmds := commands{
		command: make(map[string]func(*state, command) error),
	}

	// Check for correct number of arguments
	if len(os.Args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
	}

	// Register commands
	cmds.register("login", handlerLogins)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	err = cmds.run(currentState, command{Name: os.Args[1], Args: os.Args[2:]})

	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}

// middlware function to verify logged in user
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
