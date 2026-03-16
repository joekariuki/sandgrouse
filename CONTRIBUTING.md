# Contributing to sandgrouse

Thank you for your interest in contributing to sandgrouse. This project exists because developers on metered connections deserve the same AI tool experience as everyone else.

## How to help

### Share your experience

The most valuable contribution right now is your story. If you've experienced the bandwidth problem this project solves, open an issue and tell us:

- Where are you based?
- What AI tools do you use?
- What type of connection are you on (mobile data, metered broadband, satellite)?
- How does bandwidth cost affect your AI tool usage?

These stories shape our priorities and help us build the right thing.

### Report bugs

If something doesn't work, open an issue with:

- Your operating system and version
- How you installed sandgrouse (npm, brew, binary)
- What LLM tool you were using
- What happened vs what you expected
- Any error output from `sg status` or the terminal

### Suggest features

Open an issue with the "feature request" label. Describe:

- What problem you're trying to solve
- How you currently work around it
- What you'd like sandgrouse to do instead

### Contribute code

1. Fork the repo
2. Create a branch (`git checkout -b feature/your-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit with a clear message
6. Push and open a pull request

#### Development setup

```bash
# Prerequisites
# - Go 1.22+
# - Node.js 18+ (for npm wrapper testing)

# Clone
git clone https://github.com/jokariuki/sandgrouse.git
cd sandgrouse

# Build
go build -o sg ./cmd/sg

# Run
./sg start

# Test
go test ./...

# Lint
golangci-lint run
```

#### Code style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Write tests for new functionality
- Keep functions small and focused
- Comment exported functions and types
- No dependencies without discussion (open an issue first)

#### Areas where contributions are welcome

- Additional LLM provider support (Gemini, Mistral, Cohere, etc.)
- Compression algorithm improvements
- Dashboard UI enhancements
- Documentation and examples
- Translations of the README and manifesto
- Testing on different platforms and network conditions
- Performance benchmarking

## Code of conduct

Be kind. Be respectful. Remember that this project serves a global community of developers, many of whom are working under constraints that may be unfamiliar to you. We're all here to make AI tools more accessible.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
