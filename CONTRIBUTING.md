# Contributing to Schluckauf

Thank you for considering contributing to Schluckauf! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Bugs

If you find a bug, please [open an issue](https://github.com/fadykuzman/schluckauf/issues/new) with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Docker version, browser)
- Screenshots if applicable

### Suggesting Features

Feature requests are welcome! Please [open an issue](https://github.com/fadykuzman/schluckauf/issues/new) describing:
- The problem you're trying to solve
- Your proposed solution
- Any alternatives you've considered
- How it fits into the [roadmap](ROADMAP.md)

### Submitting Pull Requests

Pull requests are welcome! You can submit PRs directly without opening an issue first, but please:

1. **Keep PRs small and focused** - One feature or fix per PR
2. **Link to issues** - If your PR solves an existing issue, reference it
3. **Follow branch naming conventions**:
   - `feature/description` - For new features
   - `fix/description` - For bug fixes
   - `docs/description` - For documentation changes
4. **Include tests** - Add or update tests for your changes
5. **Update documentation** - Update README.md and other docs as needed

## Development Setup

### Option 1: Docker Development

The fastest way to get started:

```bash
# Clone the repository
git clone https://github.com/fadykuzman/schluckauf.git
cd schluckauf

# Create required directories
mkdir -p photos scans trash data

# Start the application
docker-compose up --build

# Open in browser
open http://localhost:8087
```

Make changes, rebuild with `docker-compose up --build`, and test.

### Option 2: Local Development

For faster iteration without Docker:

**Prerequisites:**
- [Go 1.25.1+](https://go.dev/dl/)
- [SQLite](https://sqlite.org/download.html)

**Setup:**
```bash
# Clone the repository
git clone https://github.com/fadykuzman/schluckauf.git
cd schluckauf

# Create required directories
mkdir -p photos scans trash data

# Run the application
go run cmd/dup-reviewer/main.go

# Open in browser
open http://localhost:8080
```

The frontend is served from the `web/` directory. Changes to HTML/CSS/JS are reflected immediately on refresh.

### Generating Test Data

For testing, use the included script to generate sample duplicate images:

```bash
./scripts/download-images.sh
```

**Prerequisites:** [ImageMagick](https://imagemagick.org/script/download.php) must be installed.

## Code Style Guidelines

### Go Code

- **Formatting**: Use `go fmt` to format all Go code
- **Linting**: Run `golint` before submitting PRs
- **No comment requirements**: Write comments where helpful, but not mandatory
- **Keep it simple**: Prefer standard library over external dependencies

**Run before committing:**
```bash
go fmt ./...
golint ./...
go test ./...
```

### Frontend Code (HTML/CSS/JS)

- **Formatting**: Use standard formatting conventions
- **Linting**: Run ESLint for JavaScript code
- **Consistency**: Follow existing code patterns and structure
- **No frameworks**: Keep the vanilla stack (no React, Vue, etc.)

**Run before committing:**
```bash
# If you have ESLint configured
eslint web/*.js
```

## Testing

### Manual Testing

All changes should be manually tested. Reference [MANUAL_TESTS.md](MANUAL_TESTS.md) for the core test scenarios:

1. Application startup
2. Scanning for duplicates
3. Keyboard navigation
4. Image selection and marking
5. File operations (move to trash)
6. Data persistence

### Automated Tests

If your PR adds new functionality:
- Add Go tests where applicable
- Update MANUAL_TESTS.md with new test scenarios

### Testing Checklist

Before submitting a PR:
- [ ] Code builds without errors
- [ ] All existing tests pass
- [ ] New functionality is tested (manual or automated)
- [ ] Tested in Docker environment
- [ ] No console errors in browser
- [ ] README updated if needed

## Pull Request Process

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the code style guidelines
3. **Test thoroughly** using manual tests and/or automated tests
4. **Update documentation** (README, ROADMAP, etc.) if needed
5. **Submit the pull request** with:
   - Clear title describing the change
   - Description of what changed and why
   - Reference to related issues (if any)
   - Screenshots for UI changes

**PR Review:**
- All PRs are reviewed personally by the maintainer
- Feedback will be provided within a reasonable timeframe
- Please be patient and responsive to review comments

## Project Structure

```
schluckauf/
├── cmd/
│   └── dup-reviewer/     # Application entry point
├── internal/
│   ├── handler/          # HTTP handlers and API
│   ├── loader/           # Czkawka JSON parsing
│   └── storage/          # SQLite database operations
├── web/                  # Frontend files (HTML/CSS/JS)
├── scripts/              # Utility scripts
├── MANUAL_TESTS.md       # Manual testing scenarios
├── ROADMAP.md            # Project roadmap
└── docker-compose.yml    # Docker setup

```

## Architecture Principles

- **Privacy-first**: No external API calls, all data local
- **Simplicity**: Vanilla stack, minimal dependencies
- **Speed**: Keyboard-driven workflow for fast operations
- **Safety**: Move to trash, never permanent deletion

See [CLAUDE.md](CLAUDE.md) for detailed architecture documentation.

## Questions?

If you have questions about contributing, please [open an issue](https://github.com/fadykuzman/schluckauf/issues/new) with your question. We're happy to help!

## Code of Conduct

Be respectful, constructive, and professional. We're all here to build something useful together.

## License

By contributing to Schluckauf, you agree that your contributions will be licensed under the MIT License.
