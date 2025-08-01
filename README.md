# Gator CLI

Gator is a command-line tool for managing RSS feeds and users, built in Go and backed by PostgreSQL.

## Prerequisites

- **Go** (version 1.18 or newer recommended): [Install Go](https://golang.org/doc/install)
- **PostgreSQL**: [Install PostgreSQL](https://www.postgresql.org/download/)

## Installation

Clone the repository and install the CLI using `go install`:

```sh
git clone https://github.com/nunseik/gator.git
cd gator
go install .
```

This will build and install the `gator` binary to your `$GOPATH/bin`.

## Configuration

Create a configuration file named `config.yaml` in the root directory. Example:

```yaml
DBURL: "postgres://user:password@localhost:5432/gatordb?sslmode=disable"
CurrentUserName: ""
```

Make sure your PostgreSQL instance is running and the database exists.

## Usage

Run the CLI with commands:

```sh
gator <command> [arguments...]
```

Some available commands:

- `register <username>`: Register a new user.
- `login <username>`: Log in as an existing user.
- `users`: List all users.
- `addfeed <feed_name> <feed_url>`: Add a new RSS feed (requires login).
- `feeds`: List all feeds.
- `follow <feed_url>`: Follow a feed (requires login).
- `following`: List feeds you are following (requires login).
- `unfollow <feed_url>`: Unfollow a feed (requires login).
- `browse <feed_url> [limit]`: Browse posts from a feed.

## Future Plans

- **Major Refactor:** Commands will be separated into multiple files by category (e.g., user commands, feed commands, post commands) for better maintainability and clarity.
- Additional features and improvements are planned.

---

Feel free to open issues or contribute!
