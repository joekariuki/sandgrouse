# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-03-26

### Added

- Core HTTP reverse proxy that transparently forwards LLM API traffic through localhost:8080
- Gzip and brotli compression enforcement on all requests and responses (70-80% bandwidth reduction)
- Context deduplication engine that eliminates repeated system prompts, unchanged files, and re-sent conversation history
- CLI with commands: `sg start`, `sg stop`, `sg status`, `sg stats`, `sg config`
- Daemon mode with PID file management for background operation
- Web dashboard at localhost:8585 with real-time bandwidth savings visualization
- Per-request bandwidth tracking with SQLite storage and 90-day retention
- YAML configuration at `~/.sandgrouse/config.yml` with sensible defaults
- Anthropic API provider support
- OpenAI API provider support (works with any OpenAI-compatible API)
- Google Gemini API provider support
- Distribution via npm (`npx sandgrouse`), Homebrew (`brew install sandgrouse`), and direct binary download
- Single binary with zero runtime dependencies
