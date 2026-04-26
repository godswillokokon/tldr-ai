# tldr-ai

This repository contains the **tldr-ai-be** Go service. See [`tldr-ai-be/README.md`](tldr-ai-be/README.md) for API and environment details.

## Environment files

- **`tldr-ai-be/.env.example`** — committed template with safe defaults and every supported variable. No real secrets.
- **Local `.env`** — copy from the example, add `ANTHROPIC_API_KEY` and other secrets. **Never commit** `.env`; it is ignored at the repo root and under `tldr-ai-be/` (see `.gitignore` files).

At startup, the server loads only `.env` files (repo root then `tldr-ai-be/.env`); see `tldr-ai-be/README.md`. `.env.example` is a copy template only.
