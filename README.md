# Askeladden

Askeladden is a Discord bot for Norwegian language communities, with a focus on grammar correction and community engagement.

## Features

### 🔨 Banned Word System
- **Report incorrect words**: React with 🔨 emoji to report grammatically incorrect words
- **Dual approval**: Words require approval from both Opplysar and Rettskrivar roles
- **Forum discussions**: Approved words automatically get forum threads for community discussion
- **Real-time warnings**: Bot warns users when they use banned words with links to discussions

### ❓ Question of the Day
- **Community questions**: Users can submit questions that get approved by moderators
- **Scheduled posting**: Bot posts approved questions on a schedule
- **Fair distribution**: Questions are distributed evenly to ensure all get asked

### ⭐ Starboard
- **Highlight messages**: Star messages to feature them in a dedicated starboard channel
- **Configurable threshold**: Set minimum stars required for starboard inclusion

### 🔐 Role-based Permissions
- **Granular control**: Different roles can approve different types of content
- **Combined approvals**: Some features require multiple role approvals for added quality control

## Building and Running

### Building

To build the bot, run the following command. This will create an executable named `askeladden`.

```bash
 go build ./cmd/askeladden
```

### Running

You can run the bot with different configurations by setting the `CONFIG_FILE` and `SECRETS_FILE` environment variables. If these variables are not set, the bot will default to `config.yaml` and `secrets.yaml` respectively.

**Production:**

```bash
./askeladden
```

**Development:**

To run the bot in beta mode, a handy script is provided.

## Documentation

### Discord Embed Guidelines

For standardized Discord embed creation and consistent styling, see the [Discord Embed Guidelines](docs/EMBEDS.md).

The embed system provides:
- Consistent colors and styling across all bot interactions
- Fluent builder interface for complex embeds
- Helper functions for common embed patterns
- Reduced code duplication and improved maintainability


## Production Deployment

Deploying to production on `heim.bitraf.no` is done with a single script. Here's how to do it:

### 1. Configure Production Settings

Make sure you have the following files in the root directory:

- `config.yaml`: The main configuration file. All production channel and role IDs must be set here.
- `secrets.yaml`: A file with the database password and Discord bot token:

  ```yaml
  database:
    user: <your_username>
    password: <your_password>
  discord:
    token: <your_bot_token>
  ```

### 2. Run the Deployment Script

Run the build and deployment script:

```bash
./build-and-deploy.sh
```

The script will:
- **Build** the binary for linux
```bash
GOOS=linux GOARCH=amd64 go build -o askeladden ./cmd/askeladden
```
- **Stop** any existing bot processes on the server
- **Copy** the binary, config, and secrets files
- **Start** the bot in a new `tmux` session


### 3. Verify the Deployment

The script will show you the latest output from the bot. You can also manually check the bot's status with:

```bash
# View the tmux session
ssh ellinorlinnea@heim.bitraf.no 'tmux attach -t askeladden'

# View the latest logs
ssh ellinorlinnea@heim.bitraf.no 'tmux capture-pane -t askeladden -p'
```

To detach from the tmux session, press **Ctrl+B**, then **D**.