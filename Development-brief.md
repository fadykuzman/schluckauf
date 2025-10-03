# Duplicate Image Reviewer - Development Brief

## 🎯 Project Goal
Build a **privacy-first, self-hosted web app** to review and clean up duplicate photos found by Czkawka CLI.

---

## 🏗️ Tech Stack
- **Backend**: Go (stdlib only: `net/http`, `encoding/json`, `os/exec`)
- **Database**: SQLite (pure-Go driver: `modernc.org/sqlite`)
- **Frontend**: Vanilla HTML + CSS + JavaScript (no frameworks)
- **Deployment**: Docker + docker-compose

---

## 📋 Core Features (MVP)

### 1. Load Duplicates
- Parse JSON from `czkawka_cli` output
- Store duplicate groups in SQLite

### 2. Web UI
- Show duplicate groups (one at a time)
- Display images side-by-side with metadata (resolution, size, path)
- Progress indicator: "Group X of Y"

### 3. User Actions
- **Keep** button per image (leaves original in place)
- **Trash** button per image (moves to `./trash/` directory)
- **Next Group** button
- **Keyboard shortcuts**: `1-9` select image, `K` keep, `D` trash, `N` next

### 4. State Persistence
- Save decisions to SQLite
- Resume session after restart

---

## 🔒 Privacy Requirements (GDPR)

### Must Have
- ❌ No external API calls
- ❌ No telemetry/analytics
- ✅ Read-only access to photos
- ✅ All data stored locally
- ✅ Export decisions as JSON (GET `/api/export`)
- ✅ Delete all data (DELETE `/api/data/clear`)

### File Structure
```
/photos         # User's photos (read-only, mounted volume)
/data           # SQLite database
/trash          # Deleted files go here
/scans          # Czkawka JSON output
```

---

## 📁 Project Structure

```
dup-reviewer/
├── cmd/
│   └── dup-reviewer/
│       └── main.go              # Entry point
├── internal/
│   ├── loader/
│   │   └── czkawka.go          # Parse czkawka JSON
│   ├── storage/
│   │   └── sqlite.go           # Database ops
│   └── handler/
│       ├── api.go              # HTTP API handlers
│       └── files.go            # Serve images safely
├── web/
│   ├── index.html
│   ├── app.js                  # Frontend logic
│   └── style.css
├── Dockerfile                   # Multi-stage build
├── docker-compose.yml
└── README.md
```

---

## 🗄️ Data Models

### Czkawka JSON Input
```json
{
  "duplicates": [
    {
      "size": 2156432,
      "files": [
        "/photos/IMG_1234.jpg",
        "/photos/backup/IMG_1234.jpg"
      ]
    }
  ]
}
```

### SQLite Schema
```sql
CREATE TABLE groups (
    id INTEGER PRIMARY KEY,
    hash TEXT UNIQUE,
    size INTEGER,
    file_count INTEGER
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    group_id INTEGER,
    path TEXT,
    resolution TEXT,
    filesize INTEGER,
    action TEXT DEFAULT 'pending',  -- 'keep' | 'trash' | 'pending'
    FOREIGN KEY(group_id) REFERENCES groups(id)
);

CREATE INDEX idx_group_action ON files(group_id, action);
```

---

## 🔌 API Endpoints

```
GET  /                          → Serve web UI
GET  /api/groups                → List all duplicate groups
GET  /api/groups/:id            → Get specific group with files
POST /api/groups/:id/files/:fid → Mark file as keep/trash
GET  /api/image?path=...        → Serve image (with path validation)
GET  /api/export                → Export decisions as JSON
DELETE /api/data/clear          → Delete all stored data
GET  /health                    → Health check
```

---

## 🎨 UI Layout

```
┌─────────────────────────────────────────────┐
│  Duplicate Reviewer      Group 47 / 235    │
├─────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────┐    │
│  │ Image 1 │  │ Image 2 │  │ Image 3 │    │
│  │ 4032x   │  │ 4032x   │  │ 1920x   │    │
│  │ 3024    │  │ 3024    │  │ 1080    │    │
│  │ 2.1 MB  │  │ 2.1 MB  │  │ 847 KB  │    │
│  │[Keep][❌]│  │[Keep][❌]│  │[Keep][❌]│    │
│  └─────────┘  └─────────┘  └─────────┘    │
│                                             │
│  Shortcuts: 1-9=Select  K=Keep  D=Trash   │
│             N=Next Group                    │
└─────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Development
```bash
# Run locally
go run cmd/dup-reviewer/main.go

# Open browser
open http://localhost:8080
```

### Production (Docker)
```bash
# Build and run
docker-compose up -d

# Import duplicates
docker-compose --profile scan run czkawka \
  dup -d /photos -f /scans/duplicates.json --export-json
```

---

## ⚠️ Security Requirements

### Path Validation (Critical!)
```go
// Prevent directory traversal attacks
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
- Mount photos as read-only
- No outbound network access needed

---

## 📦 Dependencies

```go
// go.mod
module github.com/yourusername/dup-reviewer

go 1.21

require (
    modernc.org/sqlite v1.28.0  // Pure-Go SQLite
)
```

---

## ✅ Acceptance Criteria

- [ ] Parse czkawka JSON and load into SQLite
- [ ] Web UI displays duplicate groups
- [ ] Can mark files as keep/trash via buttons
- [ ] Keyboard shortcuts work (1-9, K, D, N)
- [ ] Files moved to `./trash/` when trashed
- [ ] State persists across restarts
- [ ] Export decisions as JSON
- [ ] Docker deployment works
- [ ] No external network calls
- [ ] Path traversal protection

---

## 🎯 Success Metrics

**User can review 100 duplicate groups in under 5 minutes using keyboard shortcuts.**

---

## 📝 Notes for Claude Code

1. Start with backend (parse JSON + SQLite)
2. Add HTTP server with file serving
3. Build frontend with keyboard shortcuts
4. Add Docker configuration last
5. Keep it simple - no frameworks, minimal dependencies
6. Focus on speed and keyboard-driven workflow

**Priority**: Functional MVP over polish. Get keyboard shortcuts working first!
