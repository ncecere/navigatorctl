# navigatorctl

**navigator cuddle**

A command-line interface for managing teams, users, and API keys.

## Installation

### Using Go

If you have Go installed:
```bash
go install github.com/ncecere/navigatorctl@latest
```

### Pre-built Binaries

#### macOS

Using Homebrew:
```bash
# Coming soon...
# brew install ncecere/tap/navigatorctl
```

Manual installation:
```bash
# For Intel Macs (x86_64)
curl -L https://github.com/ncecere/navigatorctl/releases/latest/download/navigatorctl_Darwin_x86_64.tar.gz | tar xz
chmod +x navigatorctl
sudo mv navigatorctl /usr/local/bin/

# For Apple Silicon Macs (arm64)
curl -L https://github.com/ncecere/navigatorctl/releases/latest/download/navigatorctl_Darwin_arm64.tar.gz | tar xz
chmod +x navigatorctl
sudo mv navigatorctl /usr/local/bin/
```

#### Linux

Using script:
```bash
# For x86_64
curl -L https://github.com/ncecere/navigatorctl/releases/latest/download/navigatorctl_Linux_x86_64.tar.gz | tar xz
chmod +x navigatorctl
sudo mv navigatorctl /usr/local/bin/

# For arm64
curl -L https://github.com/ncecere/navigatorctl/releases/latest/download/navigatorctl_Linux_arm64.tar.gz | tar xz
chmod +x navigatorctl
sudo mv navigatorctl /usr/local/bin/
```

#### Windows

1. Download the latest release from [GitHub Releases](https://github.com/ncecere/navigatorctl/releases/latest)
2. Extract the ZIP file
3. Add the extracted directory to your PATH environment variable
4. Open a new PowerShell window and verify the installation:
```powershell
navigatorctl --version
```

### Verifying Installation

After installation, verify it works:
```bash
navigatorctl --version
```

## Configuration

The CLI can be configured using a configuration file or environment variables.

### Configuration File

Create a configuration file at `~/.navigatorctl.yaml`:

```yaml
api:
  url: https://ai.bitop.dev
  key: your-api-key
```

### Environment Variables

You can also use environment variables:

```bash
export NAVIGATORCTL_API_URL=https://ai.bitop.dev
export NAVIGATORCTL_API_KEY=your-api-key
```

## Usage

### Global Flags

- `--api-url`: API endpoint URL (overrides config)
- `--api-key`: API key for authentication (overrides config)
- `--output, -o`: Output format (table, json)

### Team Commands

#### List Teams
```bash
navigatorctl team list
```

#### Team Information
```bash
# Using team ID
navigatorctl team info --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671

# Using team alias
navigatorctl team info --team-alias MYTEAM
```

#### Team Members
```bash
# List members
navigatorctl team members --team-alias MYTEAM

# Add member
navigatorctl team add-member --team-alias MYTEAM --user-id user123 --email user@example.com --role user

# Remove member
navigatorctl team remove-member --team-alias MYTEAM --user-id user123
```

#### Team API Keys
```bash
navigatorctl team keys --team-alias MYTEAM
```

### User Commands

#### User Information
```bash
# Using user ID
navigatorctl user info --user-id user123

# Using email
navigatorctl user info --email user@example.com
```

#### User Teams
```bash
navigatorctl user teams --user-id user123
```

#### User API Keys
```bash
navigatorctl user keys --user-id user123
```

### Key Commands

#### List Keys
```bash
navigatorctl key list --api-url https://ai.bitop.dev --api-key sk-6425
```

#### Key Information
```bash
navigatorctl key info --key <key_string> --api-url https://ai.bitop.dev --api-key sk-6425
```

### Output Formats

All commands support both table and JSON output:

```bash
# Default table format
navigatorctl team list

# JSON format
navigatorctl team list --output json
```

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ncecere/navigatorctl.git
cd navigatorctl

# Build
go build -o navigatorctl

# Run tests
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
