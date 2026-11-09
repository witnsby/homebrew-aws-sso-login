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

The formula in `Formula/aws-sso-login.rb` is **automatically kept up to date** by a Go script (`src/main.go`) that runs via GitHub Actions. The [upstream repository](https://github.com/witnsby/aws-sso-login) starts it with the new tag as soon as it publishes a release; it also runs on every push to `main`, once a day on a schedule, and on manual dispatch:

1. Fetches the requested release (the `tag` input), or the latest release when no tag is given. The run fails when the release does not exist or lacks one of the four platform binaries.
2. Downloads each platform binary.
3. Computes the SHA256 checksum for each binary.
4. Rewrites the formula with the new version, URLs, and checksums.
5. Commits and pushes the updated formula back to `main` (skipped when the formula is already current).

You do not need to open a PR to update the formula — the CI pipeline handles it automatically. If the upstream dispatch did not run (for example, its `HOMEBREW_TAP_TOKEN` secret expired), the next scheduled run picks up the latest release; to update right away, run the workflow manually, optionally for a specific tag:

```bash
gh workflow run publish_release.yml -R witnsby/homebrew-aws-sso-login -f tag=v0.4.0
```

## Development

### Prerequisites

- Go 1.23.4+
- `make`

### Run the updater locally

```bash
make install-deps   # go mod tidy
make run            # update Formula/aws-sso-login.rb to the latest release
RELEASE_TAG=v0.3.0 make run   # ... or to a specific release
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
