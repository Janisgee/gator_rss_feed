# Gator - RSS feed arrgregator

## Introduction

Gator is a multi-user CLI application.
It is for adding RSS feeds from across the internet to be collected. And store the collected posts in a PostgreSQL database.
Follow and unfollow RSS feeds that other users have added.

## Installation

This Gator application requires Golang and Postgres installation to run the program.

#### 1. Install Go 1.22 or later

The Gator requires a Golang installation, and only works on Linux and Mac. If you're on Windows, you'll need to use WSL. Make sure you install go in your Linux/WSL terminal, not your Windows terminal/UI. There are two options:

Option 1: The webi installer is the simplest way for most people. Just run this in your terminal:

```
curl -sS https://webi.sh/golang | sh
```

Option 2: Use the official installation instructions. For further information, you can read https://go.dev/doc/install.

#### 2. Install PostgresSQL

###### Mac OS with brew

```
brew install postgresql@16
```

###### Linux / WSL (Debian). Here are the docs from Microsoft, but simply:

```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

#### 3. Install Gator

```bash
go install ...
```

## Configuration

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the value with your database connection string.

## Usage

#### Create a new user:

```bash
gator register <name>
```

#### Add a feed:

```bash
gator addfeed <url>
```

#### Start the aggregator:

```bash
gator agg 30s
```

#### View the posts:

```bash
gator browse [limit]
```

#### There are a few other commands you'll need as well:

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed that already exists in the database
