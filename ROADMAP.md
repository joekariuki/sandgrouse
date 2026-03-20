# Roadmap

## Current focus

Building Phase 1 — the core CLI proxy with compression, deduplication, dashboard, and multi-provider support. Target launch: March 26, 2026.

## Phase 1: CLI Proxy (v0.1.0)

Local compression proxy for developers using Claude Code, Cursor, and OpenAI-compatible tools on metered connections.

- Transparent HTTP proxy with gzip/brotli compression
- Context deduplication for repeated content
- CLI commands (`sg start`, `sg stop`, `sg status`, `sg config`)
- Real-time web dashboard with bandwidth savings visualization
- Bandwidth tracking with SQLite storage
- Anthropic, OpenAI, and Google Gemini provider support
- Distribution via npm, Homebrew, and direct binary download

## Phase 1.5: SDKs + API

Python and Node.js libraries that bring bandwidth optimization to any LLM-powered application, backed by a local REST API that any tool or agent can talk to.

- Local REST API on localhost:8585 for programmatic access to bandwidth stats, session data, and proxy control
- Python SDK with httpx transport wrapper and ASGI/WSGI middleware
- Node.js SDK with fetch wrapper and Express/Fastify middleware
- Drop-in `optimize()` for Anthropic and OpenAI client libraries
- Published to PyPI and npm

## Phase 1.75: MCP Server

An MCP server that lets AI agents query and control sandgrouse directly. Claude Code, Codex, and other agentic tools can check bandwidth usage, adjust compression settings, and make informed decisions about how they consume data.

- MCP server exposing sandgrouse stats, configuration, and session data as tools
- Agents can query: bytes saved, compression ratio, request history, per-provider breakdown
- Agents can act: pause proxy, change compression level, clear cache
- Enables self-aware bandwidth management (an agent can see how much data it's consuming and adjust)
- Published as an MCP server package

## Phase 2: Desktop App

A desktop application that makes sandgrouse accessible without a terminal.

- Tauri-based app with system tray integration
- Automatic LLM tool detection and proxy configuration
- Visual dashboard with enhanced analytics
- Ollama smart routing for local model traffic
- Auto-start on login
- Installers for macOS, Linux, and Windows

## Phase 3: Browser Extension

Optimize LLM web traffic for anyone, regardless of technical ability.

- Manifest V3 extension for Chrome and Firefox
- Fetch/XHR interception on Claude.ai, ChatGPT, and Gemini
- Auto-routing through the CLI proxy when running
- Popup dashboard with bandwidth savings badge

## Phase 4: Mobile

Bandwidth optimization for LLM apps on iOS and Android.

- Go mobile library for shared compression logic
- Android VpnService implementation
- iOS NetworkExtension implementation
- Native UI for both platforms

## Beyond

The long-term vision is a set of standards and protocols for bandwidth-efficient AI communication. If sandgrouse succeeds, these ideas — content deduplication, delta encoding, intelligent compression — should be adopted by AI providers themselves.

The API-first approach ensures sandgrouse is not just a tool humans use, but infrastructure that agents and applications build on. As software becomes increasingly agentic — computers talking to computers — the proxy, the API, and the MCP server together make bandwidth optimization a primitive that any system can leverage.

---

Want to help? See [CONTRIBUTING.md](CONTRIBUTING.md) or open an [issue](https://github.com/jokariuki/sandgrouse/issues).
