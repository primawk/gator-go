This project requires Postgres and Go installed.


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
./your-cli-app register alice

# Log in as an existing user
./your-cli-app login alice

# Add a new feed (requires login)
./your-cli-app addfeed "Boot.dev Blog" "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"

# View all available feeds
./your-cli-app feeds

# Follow a feed (requires login)
./your-cli-app follow "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"

# Browse the 5 most recent posts
./your-cli-app browse 5