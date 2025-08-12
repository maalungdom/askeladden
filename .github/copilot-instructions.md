# Askeladden Discord Bot

Askeladden is a Discord bot for Norwegian language communities, with a focus on grammar correction and community engagement. The bot is written in Go 1.24+ and uses DiscordGo, MySQL database, and YAML configuration.

**CRITICAL**: Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Initial Setup and Build
- **Download Go modules**: `go mod download` -- takes ~5 seconds
- **Build main application**: `go build ./cmd/askeladden` -- takes ~18 seconds first time, ~0.4 seconds subsequent builds. NEVER CANCEL.
- **Build for production (Linux)**: `GOOS=linux GOARCH=amd64 go build -o askeladden ./cmd/askeladden` -- takes ~0.4 seconds
- **Build beta version**: `go build -o askeladden-beta ./cmd/askeladden` -- takes ~0.4 seconds

### Code Quality and Validation
- **Format code**: `go fmt ./...` -- takes ~0.2 seconds. ALWAYS run before committing.
- **Vet code**: `go vet ./cmd/askeladden ./internal/...` -- takes ~2.4 seconds. Exclude `/tools` due to multiple main functions.
- **Test main packages**: `go test ./cmd/askeladden ./internal/...` -- instant (no test files exist)
- **Full test (will fail)**: `go test ./...` -- fails due to multiple main functions in `/tools/` directory

### Running the Application
- **Normal run** (will fail without config): `./askeladden`
- **With custom config**: `CONFIG_FILE=config/config-beta.yaml SECRETS_FILE=config/secrets-beta.yaml ./askeladden`
- **Beta mode**: `./run-beta.sh` -- requires `config-beta.yaml` and `secrets-beta.yaml` in root directory
- **Production deployment**: `./build-and-deploy.sh` -- builds, deploys to heim.bitraf.no server (requires SSH access)

### Testing Tools
- **Help text verification**: `./tools/test_help/test_help.go` -> build with `go build ./tools/test_help` -> run `./test_help`
- **Command matching test**: `./tools/test_match/test_match.go` -> build with `go build ./tools/test_match` -> run `./test_match` (expects runtime error without proper bot context)

## Validation

### CRITICAL Build Requirements
- **Go version**: Go 1.24+ required (check with `go version`) - tested with 1.24.5
- **NEVER CANCEL**: All builds complete in under 30 seconds. Set timeout to 60+ seconds for safety.
- **Dependencies**: All managed via `go.mod` - no external package managers needed

### Configuration Requirements
The bot requires two YAML files to run:
- `config/config.yaml` - main configuration (channel IDs, database settings)
- `config/secrets.yaml` - sensitive data (Discord token, database password)
- Only `config/config-beta.yaml` exists in the repository for reference
- Beta script `run-beta.sh` expects files named `config-beta.yaml` and `secrets-beta.yaml` in root directory

### Manual Validation Steps
- **ALWAYS run `go fmt ./...` before making any commit**
- **ALWAYS run `go vet ./cmd/askeladden ./internal/...` to check for issues**
- **Build and verify binary creation**: Check that `./askeladden` executable is created and is ~11-12MB
- **Test help text**: Build and run `./test_help` to verify command text generation (should output ~12 lines)
- **Verify production build**: Cross-compile with `GOOS=linux GOARCH=amd64` flag works correctly
- **Manual functional scenarios**: 
  - Help text contains Norwegian commands (`hei`, `hjelp`, `spør`) and emoji prefixes
  - Binary is reasonable size (~12MB indicates successful build with all dependencies)
  - Test tools build and run without compilation errors

## Common Tasks

### Working with the Codebase
- **Main entry point**: `/cmd/askeladden/main.go`
- **Core bot logic**: `/internal/bot/bot.go`
- **Discord handlers**: `/internal/bot/handlers/`
- **Business services**: `/internal/bot/services/`
- **Commands**: `/internal/commands/` - all bot commands defined here
- **Configuration**: `/internal/config/config.go`
- **Database**: `/internal/database/` - MySQL database operations
- **Discord embeds**: Follow guidelines in `/docs/EMBEDS.md`

### Key Features to Understand
- **Banned Word System**: React with 🔨 to report incorrect words
- **Question of the Day**: Users submit questions for scheduled posting
- **Starboard**: Star messages to feature them
- **Role-based Permissions**: Different roles for different approval levels

### Repository Structure
```
askeladden/
├── cmd/askeladden/          # Main application entry point
├── internal/                # Core application logic
│   ├── bot/                 # Bot core and handlers
│   ├── commands/            # Discord commands
│   ├── config/              # Configuration management
│   ├── database/            # Database operations
│   ├── permissions/         # Role-based permissions
│   └── reactions/           # Reaction handlers
├── config/                  # Configuration files
├── docs/                    # Documentation
├── tools/                   # Utility tools and scripts
├── build-and-deploy.sh      # Production deployment script
└── run-beta.sh              # Beta testing script
```

### Configuration Files
- **Production config**: `config/config.yaml` (not in repo)
- **Beta config**: `config/config-beta.yaml` (available for reference)
- **Secrets**: `config/secrets.yaml` (not in repo, contains Discord token and DB credentials)

### Time Expectations
- **Go module download**: ~5 seconds
- **Initial build**: ~18 seconds - NEVER CANCEL, always wait
- **Subsequent builds**: ~0.4 seconds  
- **Code formatting**: ~0.2 seconds
- **Code vetting**: ~2.4 seconds
- **Test validation tools**: ~0.4 seconds each to build

### Common Gotchas
- **Multiple main functions**: `/tools/` directory contains multiple utilities with main() functions - exclude from `go test ./...`
- **Missing configuration**: Bot will fail immediately if config files are missing (expected behavior)
- **Beta script expectations**: `run-beta.sh` expects config files in root directory, not `config/` subdirectory
- **Norwegian language**: Commands and documentation are in Norwegian (nynorsk)
- **No formal tests**: Repository uses manual validation tools instead of automated test suite

### Dependencies
From `go.mod`:
- `github.com/bwmarrin/discordgo v0.29.0` - Discord API
- `github.com/go-sql-driver/mysql v1.9.3` - MySQL driver  
- `github.com/google/uuid v1.6.0` - UUID generation
- `gopkg.in/yaml.v3 v3.0.1` - YAML configuration

### Environment Variables
- `CONFIG_FILE` - Path to main config file (default: `config/config.yaml`)
- `SECRETS_FILE` - Path to secrets file (default: `config/secrets.yaml`)

### Production Deployment
The `build-and-deploy.sh` script:
1. Builds Linux binary with `GOOS=linux GOARCH=amd64 go build -o askeladden ./cmd/askeladden`
2. Copies binary and config files to `ellinorlinnea@heim.bitraf.no:/home/ellinorlinnea/roersla/askeladden`
3. Starts bot in tmux session named `askeladden`
4. Logs to `askeladden.log`

### Beta Testing
The `run-beta.sh` script:
1. Backs up current config files
2. Copies beta configuration
3. Runs the bot with beta settings
4. Restores original configuration on exit

### Complete Validation Workflow
Run this complete sequence after making changes:
```bash
# Clean build from scratch
rm -f askeladden* test_*
go mod download                          # ~0.01 seconds
go build ./cmd/askeladden               # ~18 seconds first time, ~0.4 seconds subsequent
go fmt ./...                            # ~0.04 seconds
go vet ./cmd/askeladden ./internal/...  # ~0.24 seconds  
go build ./tools/test_help              # ~0.4 seconds
./test_help | wc -l                     # Should output ~12 lines
ls -lh askeladden                       # Should show ~12MB binary
```
The bot supports these commands (check `/tools/test_help` output):
- `hei` - Say hello to the bot
- `hjelp` - Show help message  
- `info` - Show bot information
- `kjeften` - Tell Askeladden to be quiet
- `ping` - Check if bot responds
- `spør` - Add a question for daily questions
- `tøm-db` - Clear database (admin only)
- `config` - Show current configuration
- `godkjenn` - Approve a question manually (admin only)
- `loggav` - Log off and shut down (admin only)
- `poke` - Trigger daily question manually (admin only)