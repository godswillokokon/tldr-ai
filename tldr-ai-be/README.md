# tldr-ai-be

HTTP API for text summarization (Anthropic), health, usage, and rate limits.

## Configuration

### Template: `.env.example` (not loaded at runtime)

The file `tldr-ai-be/.env.example` is safe to commit. It documents every variable (no real API keys). The server **does not** read it automatically — copy it to `.env` and edit there.

### Local secrets: `.env` (never commit)

1. Copy `tldr-ai-be/.env.example` to `tldr-ai-be/.env`.
2. Set `ANTHROPIC_API_KEY` and any other overrides.
3. Keep `tldr-ai-be/.env` out of git (see `tldr-ai-be/.gitignore`).

You can also place a repo-root `.env` for tools that run from the monorepo root; that file is listed in the root `.gitignore` as well.

### Load order (process startup)

On startup, the binary loads only (missing files are ignored):

1. `LoadDotEnvOverride(".env")` — repo root, when present
2. `LoadDotEnvOverride("tldr-ai-be/.env")` — service directory, when present

Later files override earlier ones for the same variable. Environment variables already set in the shell **are overwritten** by values in these files (override mode).

### Variables (see `.env.example`)

| Area | Variables |
|------|-----------|
| Server | `PORT` |
| Proxy / CORS | `TRUST_PROXY`, `CORS_ALLOW_ORIGIN` |
| Anthropic | `ANTHROPIC_API_KEY`, `ANTHROPIC_MODEL` |
| Rate limit | `RATE_LIMIT_RPS`, `RATE_LIMIT_MAX_IPS` |
| Usage | `USAGE_BUDGET_USD`, `USAGE_PER_CALL_USD`, `USAGE_MAX_CALLS`, `USAGE_RESET_SECRET` |

If `ANTHROPIC_API_KEY` is missing or looks like a placeholder, a one-line startup hint is logged (see `config.LogStartupEnvHint`).

## Run

```bash
cd tldr-ai-be
go run ./cmd/tldr-ai-be
```

Or with env in the shell only (no `.env` files):

```bash
PORT=3000 ANTHROPIC_API_KEY=sk-ant-... go run ./cmd/tldr-ai-be
```

Health:

```bash
curl -s http://localhost:8080/health
```

Usage snapshot:

```bash
curl -s http://localhost:8080/api/usage
```
