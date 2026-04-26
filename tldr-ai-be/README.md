# tldr-ai-be

Minimal HTTP API bootstrap.

## Run

```bash
cd tldr-ai-be
go run ./cmd/tldr-ai-be
```

Optional: set `PORT` (default `8080`).

```bash
PORT=3000 go run ./cmd/tldr-ai-be
```

Health check:

```bash
curl -s http://localhost:8080/health
```
