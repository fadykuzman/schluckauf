# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Schluckauf** (German for "hiccup") is a privacy-first, self-hosted web application for reviewing and managing duplicate photos found by Czkawka CLI. The goal is to provide a fast, keyboard-driven interface to review duplicate groups and decide which files to keep or delete.

**Core Principles:**

- Privacy-first: No external API calls, no telemetry, all data stored locally
- Simplicity: Vanilla stack (Go stdlib, vanilla HTML/CSS/JS, SQLite)
- Speed: Keyboard-driven workflow to review 100+ groups in under 5 minutes

## Tech Stack

- **Backend**: Go (stdlib only: `net/http`, `encoding/json`, `os/exec`)
- **Database**: SQLite with pure-Go driver (`modernc.org/sqlite`)
- **Frontend**: Vanilla HTML + CSS + JavaScript (no frameworks)
- **Deployment**: Docker + docker-compose

## Architecture

### High-Level Flow

1. Czkawka CLI scans photos and outputs JSON with duplicate groups
2. Loader parses JSON and stores groups/files in SQLite
3. HTTP server serves web UI and API endpoints
4. User reviews duplicates, marks files as keep/trash
5. Trashed files are moved to `./trash/` directory (not deleted)
6. State persists in SQLite for session resumption

### Module Structure

```
cmd/dup-reviewer/main.go       # Entry point, HTTP server setup
internal/loader/czkawka.go     # Parse czkawka JSON output
internal/storage/sqlite.go     # Database operations
internal/handler/api.go        # HTTP API handlers
internal/handler/files.go      # Secure image serving with path validation
web/                           # Static frontend files
```

### Database Schema

```sql
groups: id, hash (unique), size, file_count
files: id, group_id, path, resolution, filesize, action ('keep'|'trash'|'pending')
```

## API Endpoints

- `GET /` - Serve web UI
- `GET /api/groups` - List all duplicate groups
- `GET /api/groups/:id` - Get specific group with files
- `POST /api/groups/:id/files/:fid` - Mark file action (keep/trash)
- `GET /api/image?path=...` - Serve image (with path validation)
- `GET /api/export` - Export decisions as JSON
- `DELETE /api/data/clear` - Delete all stored data
- `GET /health` - Health check

## Development Commands

### Running locally

```bash
go run cmd/dup-reviewer/main.go
# Server runs on http://localhost:8080
```

### Docker deployment

```bash
docker-compose up -d

# Import duplicates from czkawka
docker-compose --profile scan run czkawka \
  dup -d /photos -f /scans/duplicates.json --export-json
```

## Critical Security Requirements

### Path Traversal Protection

All image serving MUST validate paths to prevent directory traversal attacks:

```go
func validateImagePath(requestedPath, baseDir string) error {
    clean := filepath.Clean(requestedPath)
    abs, _ := filepath.Abs(clean)
    base, _ := filepath.Abs(baseDir)

    if !strings.HasPrefix(abs, base) {
        return errors.New("invalid path")
    }
    return nil
}
```

### Docker Security

- Run as non-root user (UID 1000)
- Mount `/photos` as read-only
- No outbound network access needed

### Privacy (GDPR Compliance)

- ❌ No external API calls
- ❌ No telemetry/analytics
- ✅ Read-only access to photos
- ✅ All data stored locally
- ✅ Export decisions as JSON
- ✅ Delete all data endpoint

## UI/UX Requirements

### Keyboard Shortcuts (Priority!)

- `1-9`: Select image in group
- `K`: Keep selected image
- `D`: Trash selected image
- `N`: Next group

The keyboard-driven workflow is essential for speed - focus on making this seamless.

### Display Requirements

- Show duplicate groups one at a time
- Display images side-by-side with metadata (resolution, size, path)
- Progress indicator: "Group X of Y"
- Keep/Trash buttons per image

## File System Structure

```
/photos   # User's photos (read-only, mounted volume)
/data     # SQLite database
/trash    # Deleted files moved here (not permanently deleted)
/scans    # Czkawka JSON output files
```

## Development Notes

1. **Start with backend**: Parse JSON + SQLite operations
2. **Add HTTP server**: File serving with security validation
3. **Build frontend**: Focus on keyboard shortcuts first
4. **Docker configuration**: Last step

**Priority**: Functional MVP over polish. Get the keyboard-driven workflow working before adding features.

- fit the instructions to the screen hight
- don't provide code. first guidance

