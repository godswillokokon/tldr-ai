# tldr-ai-be

HTTP API for text summarization (Anthropic), health, usage, and rate limits.

## Configuration

### Committed defaults: `.env.example`

The file `tldr-ai-be/.env.example` is safe to commit. It documents every variable (no real API keys). Values there are non-secret defaults such as `PORT=8080`, rate limit and usage defaults, and empty `ANTHROPIC_API_KEY=`.

### Local secrets: `.env` (never commit)

1. Copy `tldr-ai-be/.env.example` to `tldr-ai-be/.env`.
2. Set `ANTHROPIC_API_KEY` and any other overrides.
3. Keep `tldr-ai-be/.env` out of git (see `tldr-ai-be/.gitignore`).

You can also place a repo-root `.env` for tools that run from the monorepo root; that file is listed in the root `.gitignore` as well.

### Load order (process startup)

On startup, the binary loads, in order:

1. `LoadDotEnv(".env.example")`
2. `LoadDotEnv("tldr-ai-be/.env.example")`
3. `LoadDotEnvOverride(".env")`
4. `LoadDotEnvOverride("tldr-ai-be/.env")`

`LoadDotEnv` only fills keys that are **not** already set in the environment (so CI and shell exports win). `LoadDotEnvOverride` applies your local `.env` files and overrides. Missing files are ignored.

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
