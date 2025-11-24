# gator
Welcome to gator, a RSS feed aggregator written in Go.

Requirements:
 - PostgreSQL: https://www.postgresql.org/
 - Go https://go.dev/doc/install
   
Run command: go install "github.com/atmetz/rssgator@latest"

In the home directory, create file .gatorconfig.json containing:
{"db_url":"<postgres user name>://<postgres password>:postgres@localhost:5432/gator?sslmode=disable","CurrentUserName":""}

Command ussage:
rssgator <command> [args...]

Available commands:
login <user name>: Sets current user to new user.
register <user name> : Register a new user. Sets current user to new user.
reset : Reset database.
users : Print a list of all users.
agg <time duration> : Automatically scrapes feeds for new posts. Repeats based on time duration (etc, 1s, 1m, 1h).
addfeed <feed name> <feed url> : Add a new feed. Current user automatically follows the feed.
feeds : Print feeds followed by all current user.
follow <feed url> : Add <feed url> to current users followed feeds.
following : Lists all feeds current user is following.
unfollow <feed url> : Unfollow feed at <feed url>.
browse <limit> : Display <limit> number of latest posts in followed feeds.
