# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-30

### Added
- CLI task management application with JSON file-based storage
- `add` command to create new tasks with title and description
- `view` command to display all active tasks
- `change` command to modify existing task title and/or description
- `delete` command with soft delete mechanism (30-day retention)
- `restore` command to recover individual deleted tasks
- `restore --all` flag to recover all deleted tasks at once
- Automatic cleanup of deleted tasks older than 30 days
- Comprehensive test suite for all core functionality
- Command-line flag support with shorthand options
- Error handling and validation for all operations
- Professional README with usage examples and documentation

### Changed
- Migrated from `json.Marshal`/`json.Unmarshal` to `json.NewEncoder`/`json.NewDecoder` for better performance with large datasets
- Refactored I/O functions to use proper `os` and `io` libraries
- Improved error handling across all commands

### Fixed
- Resolved bug causing duplication of soft deleted tasks
- Fixed view command display issues

[1.0.0]: https://github.com/tristnaja/taski/releases/tag/v1.0.0
