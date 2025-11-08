# Roadmap

This document outlines the planned features and development direction for Schluckauf.

## Vision

Schluckauf aims to be a comprehensive web interface for Czkawka, providing an intuitive, keyboard-driven workflow for all file management operations. The development follows a phased approach, prioritizing media file management first, then expanding to general file operations.

## Current Status

- ✅ **Similar images detection** - Web UI for reviewing and managing similar/duplicate images found by Czkawka
- ✅ **Keyboard-driven workflow** - Fast navigation and decision-making with vim-style shortcuts
- ✅ **Safe file operations** - Move files to trash (not permanent deletion)
- ✅ **Data persistence** - SQLite-backed session management
- ✅ **Progress tracking** - Visual indicators for navigation and review progress

## Phase 1: Media Files (Current Focus)

Priority focus on media file management with visual review capabilities.

### Completed
- ✅ Similar images detection (`czkawka image`)

### Planned
1. **Similar videos detection** - Review and manage similar/duplicate videos
   - Video thumbnail previews
   - Metadata comparison (resolution, duration, codec)
   - Same keyboard-driven workflow

2. **Duplicate music files** - Manage duplicate music by tags
   - Audio metadata display (artist, album, bitrate)
   - Tag-based similarity detection
   - Batch operations for music libraries

## Phase 2: File Management

Expand to general file operations beyond media files.

### Planned
1. **Duplicate files** (`czkawka dup`)
   - Support for all file types (documents, archives, etc.)
   - Hash-based duplicate detection
   - Same review workflow as images

2. **Broken files detection** (`czkawka broken`)
   - Identify corrupted or broken files
   - File type verification
   - Batch cleanup operations

3. **Big files finder** (`czkawka big`)
   - Locate large files consuming disk space
   - Size-based filtering and sorting
   - Visual disk usage insights

## Phase 3: Cleanup & Maintenance

Complete coverage of all Czkawka operations for comprehensive file system maintenance.

### Planned
1. **Empty folders** (`czkawka empty-folders`)
   - Find and remove empty directories
   - Recursive scanning options

2. **Temporary files** (`czkawka temp`)
   - Identify temporary and cache files
   - Safe cleanup with exclusion rules

3. **Empty files** (`czkawka empty-files`)
   - Locate zero-byte files
   - Batch deletion capabilities

4. **Invalid symlinks** (`czkawka symlinks`)
   - Find broken symbolic links
   - Link target verification

5. **Invalid extensions** (`czkawka ext`)
   - Detect files with incorrect extensions
   - File type verification based on content

## Future Considerations

- **Multi-user support** - User authentication and session management
- **Scheduled scans** - Automated periodic scanning
- **Cloud storage integration** - Support for remote storage backends
- **Advanced filtering** - Custom rules and filters for file operations
- **Batch operations** - Enhanced bulk file management
- **Export/Import** - Share scan results and decisions

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to contribute to these roadmap items.

## Timeline

This is an open-source project developed in spare time. Features will be implemented based on:
- Community needs and feedback
- Contributor availability
- Technical dependencies

No strict timeline is provided, but Phase 1 is the current priority.
