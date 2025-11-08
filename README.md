# Schluckauf

Self-hosted Web UI for reviewing duplicate photos found by [Czkawka](https://github.com/qarmin/czkawka) with a speed-focused keyboard workflow.

> **Schluckauf** (shlook•owf, IPA: [ˈʃlʊkaʊ̯f]), "hiccup" in German — a playful nod to **Czkawka** (tch•kav•ka, IPA: [ˈʧ̑kafka]), "hiccup" in Polish.

## Overview

Czkawka is an excellent CLI tool for finding duplicate files with speed and accuracy. Schluckauf complements it by adding a web interface for visual review and remote access, making it easy to manage duplicate photos across local machines and remote servers.

Perfect for users managing photo libraries who want Czkawka's powerful scanning combined with the convenience of browser-based review. Deploy once with Docker, access from anywhere, and leverage keyboard shortcuts for lightning-fast decisions.

## Key Features

- **Web UI for Czkawka** — Visual interface for reviewing duplicate scan results with keyboard-driven workflow
- **Complete Backend** — HTTP API server with SQLite storage to manage scan results and user decisions
- **One-Command Deployment** — Spin up with `docker-compose up`, Czkawka bundled in the image
- **Self-Hosted & Private** — All data stays on your machine, no external API calls or telemetry
- **Keyboard-First Workflow** — Review 100+ duplicate groups in minutes with vim-style shortcuts
- **Safe Operations** — Files moved to trash directory (not deleted), easy to restore if needed

## Quick Start (POC)

> **Note:** This is a basic setup for POC testing. Production deployment will be streamlined in future releases.

### Prerequisites
- Docker and Docker Compose installed

### Setup

1. **Create required directories:**
   ```bash
   mkdir -p photos scans trash data
   ```

2. **Update docker-compose.yml** to mount your photo directory:
   ```yaml
   volumes:
     - /path/to/your/photos:/photos  # Change this to your actual photos location
   ```

3. **(Optional) Change port** in docker-compose.yml (default: 8087)

4. **Start the application:**
   ```bash
   docker-compose up --build
   ```

5. **Open in browser:** http://localhost:8087

6. **Scan for duplicates:** Enter `/photos` in the scan form and click "Scan for Duplicates"

## Important Notes

> **⚠️ Scanning Limitations:** Only scan directories that are mounted in the Docker container. Scanning other paths (e.g., `/usr`, `/etc`) will find duplicates but images cannot be displayed in the web UI. Always scan mounted volumes like `/photos`.
>
> *This limitation will be addressed post-POC with proper path validation and security restrictions.*

> **⚠️ Rescanning Behavior:** Rescanning a directory will **DELETE ALL previous scan data** including:
> - All duplicate groups (reviewed and pending)
> - Keep/Trash decisions you've made
> - Group history and metadata
>
> **Recommendation:** Complete your current review session and use "Move to Trash" to save your decisions before rescanning.
>
> **Note:** Files already moved to `/trash` are safe, but the link to which scan found them will be lost.
>
> *Post-POC: Scan history and decision preservation will be implemented.*

## Keyboard Shortcuts

### Navigation
| Keys | Action |
|------|--------|
| `↑` `↓` or `J` `K` | Navigate between groups |
| `←` `→` or `H` `L` | Navigate images within group |
| `1-9` | Jump to specific image in group |

### Actions
| Key | Action |
|-----|--------|
| `Enter` | Mark selected image as Keep |
| `Space` | Mark selected image as Trash |
| `Esc` | Deselect current image |

### Help
| Key | Action |
|-----|--------|
| `?` | Show help modal |

## Development

### Generating Test Images

For testing purposes, you can use the included script to download sample duplicate images:

```bash
# Download with defaults (70 unique images, 6 copies each)
./scripts/download-images.sh

# Custom configuration
./scripts/download-images.sh --unique=50 --copies=3

# Show help
./scripts/download-images.sh --help
```

**Options:**
- `--unique NUM` - Number of unique images to download (default: 70)
- `--copies NUM` - Number of copies per image (default: 6, use 0 for no copies)

**Prerequisites:**
- [ImageMagick](https://imagemagick.org/script/download.php) must be installed (used for image validation)

**Known Issues:**
- Some Picsum image IDs (105, 138, 148, 150, and possibly others) are unavailable or corrupted on picsum.photos
- The script will automatically retry and skip failed downloads

## Technical Stack

- **Backend:** Go 1.25.1 (no web frameworks, standard library HTTP server)
- **Database:** SQLite with WAL mode
- **Frontend:** Vanilla HTML, CSS, and JavaScript (no frameworks)
- **CLI Integration:** Czkawka v10.0.0 (bundled in Docker image)
- **Deployment:** Docker + Docker Compose

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and development phases.

**Current focus:** Phase 1 - Media file management (similar videos, duplicate music)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- How to report bugs and request features
- Development setup instructions
- Code style guidelines
- Pull request process

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Credits

This project is built as a companion tool for [Czkawka](https://github.com/qarmin/czkawka) by [qarmin](https://github.com/qarmin). Czkawka is an excellent duplicate finder licensed under MIT.
