# Technical Findings

Discoveries from building and testing the sandgrouse proxy against real LLM APIs.

---

## Finding 1: LLM APIs reject compressed request bodies

**Date:** March 20, 2026
**Updated:** March 21, 2026 (full test matrix completed)
**Affects:** Request compression (gzip/brotli on outbound request bodies)

### Summary

Neither Anthropic nor OpenAI accept compressed request bodies, despite this being a standard HTTP feature (`Content-Encoding` per RFC 9110 §8.4.1).

### Full test matrix (March 21)

| Algorithm | Header | Provider | Result | Error |
|-----------|--------|----------|--------|-------|
| brotli | Content-Encoding: br | Anthropic | Rejected | "str is not valid UTF-8: surrogates not allowed" |
| gzip | Content-Encoding: gzip | Anthropic | Rejected | "str is not valid UTF-8: surrogates not allowed" |
| gzip | Transfer-Encoding: gzip | Anthropic | Rejected | "str is not valid UTF-8: surrogates not allowed" |
| gzip | Content-Encoding: gzip | OpenAI | Rejected | "We could not parse the JSON body of your request" |

### Root cause

Both APIs parse the raw request bytes directly as JSON without checking the `Content-Encoding` header. The binary output of gzip/brotli is not valid UTF-8, so JSON parsing fails immediately at byte 0. This is a server-side implementation gap — they don't implement Content-Encoding decompression on their request ingress.

### Decision

Set `CompressRequests: false` for all providers. The compression infrastructure remains in the codebase (tested and working). If any provider adds support in the future, enabling it is a one-line change in `internal/proxy/provider.go`.

### Impact

Request compression was the highest-expected-savings layer (70-80% reduction on JSON). With it blocked, v0.1 savings come from:

1. ~~Request body compression~~ — blocked by API limitations
2. **Response compression** — works via `Accept-Encoding: gzip, br` (~37% per response)
3. **Request coalescing** — deduplicate Claude Code's double-send pattern (~50% request reduction)
4. **Context deduplication** (v0.2) — eliminates repeated content before transmission
5. **Delta encoding** (v0.3) — sends only changes between requests

### Related GitHub issues

These issues in the Claude Code repository document the same bandwidth problem:
- Issue #13911: Feature request to upgrade HTTP compression (closed as stale)
- Issue #30688: 35GB+ of outbound traffic per day
- Issue #24147: Cache read tokens consuming 99.93% of usage

---

## Finding 2: Requests are 99% of session bandwidth

**Date:** March 20, 2026
**Updated:** March 21, 2026 (full dogfood session data)

### Summary

In real Claude Code sessions, request bodies dominate bandwidth. Responses are tiny by comparison. Request sizes grow throughout a session as conversation context accumulates.

### Evidence

**Dogfood session 1** (codebase review, ~4 minutes):

| Metric | Value |
|--------|-------|
| Total requests | 53 |
| Total request data | ~6-7 MB |
| Total response data (original) | 80.9 KB |
| Response compression savings | 17.4 KB (22%) |
| Request size range | 323 B → 162 KB |

**Dogfood session 2** (adding Gemini provider, ~6 minutes):

| Metric | Value |
|--------|-------|
| Total requests | 47 |
| Total request data | ~5.5 MB |
| Total response data (original) | 79.9 KB |
| Response compression savings | 15.1 KB (19%) |
| Request size range | 323 B → 148 KB |

### Why requests are so large

1. **Full context re-transmission** — every request includes the entire conversation history. As the session grows, each request balloons from 64KB to 148KB+.
2. **File content duplication** — files read by Claude Code are included in every subsequent request, even when unchanged.
3. **System prompt repetition** — the identical system prompt (~10-15 KB) is sent on every request.
4. **Double-send pattern** — Claude Code sends each request twice (preflight + streaming), doubling request bandwidth.

### Implications

- Response compression (~37% per response) is real but affects only ~1% of total bandwidth
- The dashboard (bandwidth visibility) is the most valuable v0.1 feature
- Request coalescing (eliminating the double-send) is a quick win (~50% request reduction)
- Full context deduplication requires either API cooperation or intelligent context truncation (v0.2+)
- These numbers validate the PRD thesis and are compelling content for launch posts
