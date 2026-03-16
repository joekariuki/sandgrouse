# sandgrouse

**Stop burning data bundles on AI tools.**

Sandgrouse compresses LLM API traffic so developers on metered connections get full AI power at a fraction of the data cost.

Works with Claude Code, Codex, Cursor, Gemini, and any OpenAI-compatible tool.

```bash
npx sandgrouse    # install
sg start          # start the proxy
sg status         # see your savings
```

> A 1.2GB data bundle shouldn't disappear in two Claude Code sessions.

## What this does

Sandgrouse is a local proxy that sits between your AI tools and the cloud. It optimizes every byte before it leaves your device:

- **Compression** - Enforces gzip/brotli on all LLM API traffic (70-80% reduction)
- **Deduplication** - Eliminates repeated content: system prompts, unchanged files, re-sent conversation history
- **Delta encoding** - For files your coding tools re-read, sends only what changed
- **Dashboard** - Real-time view of bandwidth savings at localhost:8585

Everything runs locally on your machine. No data is sent anywhere except the original API destination.

## Why this exists

AI tools assume unlimited bandwidth. That assumption excludes most of the world.

If you've ever watched a data bundle vanish during a Claude Code session, rationed your AI usage based on mobile data budget, or kept working through a power outage on a phone hotspot - this is for you.

Read the full story: [MANIFESTO.md](MANIFESTO.md)

## Quick start

```bash
# Install (Node.js)
npm install -g sandgrouse

# Or via Homebrew (macOS/Linux)
brew install sandgrouse

# Or direct download
curl -fsSL https://sandgrouse.dev/install.sh | sh
```

**Set up your AI tools to use the proxy:**

```bash
# Claude Code
export ANTHROPIC_BASE_URL=http://localhost:8080

# Cursor / OpenAI-compatible tools
export OPENAI_BASE_URL=http://localhost:8080
```

**Start saving bandwidth:**

```bash
sg start          # Start the proxy
sg status         # Check savings
sg stats today    # Today's detailed stats
sg stop           # Stop the proxy
```

## How it works

```
Your AI tool (Claude Code, Cursor, etc.)
         |
         |  Uncompressed request (~250KB)
         v
   Sandgrouse proxy (localhost:8080)
         |
         |  Compressed + deduplicated (~40KB)
         v
   Cloud API (Anthropic, OpenAI)
```

Sandgrouse intercepts API requests, compresses the JSON payload, identifies and removes duplicate content (system prompts sent on every request, files that haven't changed), and forwards the optimized request. Responses are handled the same way in reverse.

Your AI tools don't know the proxy is there. The experience is identical, just cheaper.

## Configuration

Configuration lives at `~/.sandgrouse/config.yml`:

```yaml
proxy:
  port: 8080
  dashboard_port: 8585

compression:
  algorithm: brotli # brotli or gzip

providers:
  anthropic:
    enabled: true
  openai:
    enabled: true
  gemini:
    enabled: true
```

## Roadmap

- [x] Project manifesto and vision
- [ ] Core proxy with compression (v0.1)
- [ ] Context deduplication
- [ ] Web dashboard
- [ ] Python SDK (`pip install sandgrouse`)
- [ ] Node.js SDK
- [ ] Desktop app with system tray
- [ ] Browser extension for Claude.ai / ChatGPT
- [ ] Mobile optimization

See [ROADMAP.md](ROADMAP.md) for the full roadmap and [CHANGELOG.md](CHANGELOG.md) for release history.

## Why "sandgrouse"?

The sandgrouse is an African bird with a remarkable adaptation: its breast feathers absorb water like a sponge. Every morning, it flies up to 30km across the desert to a waterhole, soaks its feathers, and flies back to its chicks. The most efficient water transport system in nature.

Efficient data transport across bandwidth-scarce environments. That's what this does.

## Contributing

This project is just getting started. If you experience the bandwidth problem this solves, I especially want to hear from you.

- Star the repo if this resonates
- Open an issue to share your experience or suggest features
- PRs welcome - see [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT

---

_Built by [Joe Kariuki](https://github.com/joekariuki). Open source._
