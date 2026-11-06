# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, add RSS feeds, follow feeds, and browse posts from their followed feeds.

## Requirements

Before installing Gator, make sure you have the following installed:

- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install Gator using `go install`:

```bash
go install github.com/ahmed/gator@latest
```

Make sure your Go bin directory is included in your system's PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

You can verify the installation by running:

```bash
gator
```

## Configuration

Gator requires a configuration file named `.gatorconfig.json` in your home directory.

Create the file:

```bash
touch ~/.gatorconfig.json
```

Add the following configuration:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Update the database connection string according to your PostgreSQL setup.

## Database Setup

Start PostgreSQL:

```bash
sudo service postgresql start
```

Create the database:

```bash
sudo -u postgres createdb gator
```

Run the database migrations from the project directory:

```bash
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

## Usage

Gator is a command-line application.

### Register a user

```bash
gator register <username>
```

### Login

```bash
gator login <username>
```

### List users

```bash
gator users
```

### Add an RSS feed

```bash
gator addfeed <feed_name> <feed_url>
```

### List available feeds

```bash
gator feeds
```

### Follow a feed

```bash
gator follow <feed_url>
```

### List followed feeds

```bash
gator following
```

### Unfollow a feed

```bash
gator unfollow <feed_url>
```

### Start the feed aggregator

```bash
gator agg 1m
```

This fetches RSS feeds every minute.

### Browse posts

```bash
gator browse
```

You can also specify the number of posts:

```bash
gator browse 5
```

## Development

To run the application directly from the source code during development:

```bash
go run .
```

To build the application:

```bash
go build
```

To install the compiled CLI:

```bash
go install
```

After installation, you can run the `gator` command without using `go run`.

## Technologies

- Go
- PostgreSQL
- SQL
- SQLC
- Goose
- RSS