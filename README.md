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
go build -o askeladden
```

### Running

You can run the bot with different configurations by setting the `CONFIG_FILE` and `SECRETS_FILE` environment variables. If these variables are not set, the bot will default to `config.yaml` and `secrets.yaml` respectively.

**Production:**

```bash
./askeladden
```

**Beta:**

To run the bot in beta mode, follow these steps:

1. Ensure you are in the `/Users/eg/rørsla/askeladden` directory.

2. Use the following command to build the bot:
   ```bash
   go build ./cmd/askeladden
   ```

3. Run the bot with:


## Production Deployment

Deploying to production on `heim.bitraf.no` is done with a single script. Here's how to do it:

### 1. Build the Linux Binary

First, build the Linux binary from the `cmd/askeladden` directory:

```bash
GOOS=linux GOARCH=amd64 go build -o askeladden-linux cmd/askeladden/scheduler.go cmd/askeladden/main.go
```

### 2. Configure Production Settings

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

### 3. Run the Deployment Script

Finally, run the deployment script:

```bash
./deploy-production.sh
```

The script will:
- **Stop** any existing bot processes on the server
- **Copy** the binary, config, and secrets files
- **Start** the bot in a new `tmux` session

### 4. Verify the Deployment

The script will show you the latest output from the bot. You can also manually check the bot's status with:

```bash
# View the tmux session
ssh ellinorlinnea@heim.bitraf.no 'tmux attach -t askeladden'

# View the latest logs
ssh ellinorlinnea@heim.bitraf.no 'tmux capture-pane -t askeladden -p'
```

To detach from the tmux session, press **Ctrl+B**, then **D**.

