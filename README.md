# Gator 🐊

We're going to build an RSS feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

---

## Commands

Here is a list of available commands you can run with the CLI application:

| Command | Arguments | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `login` | `<username>` | Log in as an existing user. | No |
| `register` | `<username>` | Create a new user account and log in. | No |
| `reset` | *None* | Reset the database (clears data). | No |
| `users` | *None* | List all registered users in the database. | No |
| `agg` | `<time>` (e.g., `1s`) | Start background feed aggregation loop. | No |
| `addfeed` | `<name> <url>` | Add a new RSS feed to track. | **Yes** |
| `feeds` | *None* | List all feeds currently tracked in the database. | No |
| `follow` | `<url>` | Follow an existing feed by its URL. | **Yes** |
| `following` | *None* | List all feeds the currently logged-in user is following. | **Yes** |
| `unfollow` | `<url>` | Unfollow a feed by its URL. | **Yes** |
| `browse` | `[limit]` (optional) | View recent posts from feeds you follow. | No |

### Usage Examples

```bash
# Register a new user
./gator register alice

# Log in as an existing user
./gator login alice

# Add a new feed (requires login)
./gator addfeed "Boot.dev Blog" "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"

# View all available feeds
./gator feeds

# Follow a feed (requires login)
./gator follow "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"

# Browse the 5 most recent posts
./gator browse 5