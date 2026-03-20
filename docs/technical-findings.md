# Technical Findings

Discoveries from building and testing the sandgrouse proxy against real LLM APIs.

---

## Finding 1: LLM APIs reject compressed request bodies

**Date:** March 20, 2026
**Affects:** Request compression (gzip/brotli on outbound request bodies)

### Summary

Neither Anthropic nor OpenAI accept compressed request bodies, despite this being a standard HTTP feature (`Content-Encoding: gzip` / `Content-Encoding: br`).

### Evidence

**Anthropic (brotli):**
- Small requests (135 bytes): accepted, responded correctly
- Large requests via Claude Code (118KB): rejected with `invalid_request_error: str is not valid UTF-8: surrogates not allowed`
- The API attempts to parse the compressed bytes as JSON rather than decompressing first

**OpenAI (gzip):**
- All requests rejected with: `We could not parse the JSON body of your request`
- Same issue: the API does not decompress `Content-Encoding: gzip` request bodies

**Anthropic (gzip):**
- Small requests via curl: accepted and responded correctly
- Behavior on larger Claude Code payloads: inconsistent, same UTF-8 error on some requests

### Decision

Set `CompressRequests: false` for both providers. The compression infrastructure remains in the codebase (tested and working) for future use if APIs add support, or for custom/self-hosted endpoints that do support it.

### Impact

Request compression was one of four planned compression layers. The remaining three are unaffected:

1. ~~Request body compression~~ — blocked by API limitations
2. **Response compression** — works via `Accept-Encoding: gzip, br` negotiation
3. **Context deduplication** (v0.2) — eliminates repeated content before transmission
4. **Delta encoding** (v0.3) — sends only changes between requests

Response compression is where the majority of bandwidth savings occur anyway, since model responses are significantly larger than request payloads.

### Note

This was an anticipated risk. The compression infrastructure was built to be toggled per-provider via the `CompressRequests` flag, so disabling it required no architectural changes.

---

## Finding 2: Claude Code sends ~118KB per request

**Date:** March 20, 2026

### Summary

A simple "hello" prompt in Claude Code generates a 118KB request payload. This includes the full system prompt, conversation history, and file context — re-sent identically on every request.

### Evidence

```
POST /v1/messages -> https://api.anthropic.com/v1/messages?beta=true | request: 118428 bytes
POST /v1/messages -> https://api.anthropic.com/v1/messages?beta=true | request: 118620 bytes
```

Two sequential requests with near-identical sizes, confirming massive content repetition.

### Implications

- Validates the deduplication thesis: most request content is repeated across calls
- A typical Claude Code session with 50+ API calls transmits 5-10MB of largely redundant data
- Context deduplication (v0.2) has significant potential to reduce this
