# homebrew-aws-sso-login

Homebrew tap for [`aws-sso-login`](https://github.com/witnsby/aws-sso-login) — a CLI that streamlines AWS SSO authentication and credentials management.

## Installation

```bash
brew tap witnsby/aws-sso-login
brew install aws-sso-login
```

## Upgrade

```bash
brew upgrade aws-sso-login
```

## Uninstall

```bash
brew uninstall aws-sso-login
brew untap witnsby/aws-sso-login
```

## Supported Platforms

| OS    | Architecture |
|-------|-------------|
| macOS | x86_64 (Intel) |
| macOS | arm64 (Apple Silicon) |
| Linux | x86_64 |
| Linux | arm64 |

## How It Works

The formula in `Formula/aws-sso-login.rb` is **automatically kept up to date** by a Go script (`src/main.go`) that runs on every push to `main` via GitHub Actions:

1. Fetches the latest release from the [upstream repository](https://github.com/witnsby/aws-sso-login).
2. Downloads each platform binary.
3. Computes the SHA256 checksum for each binary.
4. Rewrites the formula with the new version, URLs, and checksums.
5. Commits and pushes the updated formula back to `main`.

You do not need to open a PR to update the formula — the CI pipeline handles it automatically.

## Development

### Prerequisites

- Go 1.23.4+
- `make`

### Run the updater locally

```bash
make install-deps   # go mod tidy
make run            # update Formula/aws-sso-login.rb to the latest release
```

### Run tests

```bash
make tests   # go test ./...
make cover   # tests with HTML coverage report
```

### Build the updater binary

```bash
make build   # output: ./tmp/get-release
```

## License

[Apache-2.0](LICENSE)
