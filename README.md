# Askeladden

Askeladden is a bot.

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
   ```bash
   CONFIG_FILE=config-beta.yaml SECRETS_FILE=secrets-beta.yaml ./askeladden-beta
   ```

