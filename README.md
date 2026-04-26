# tldr-ai

Monorepo for a small **TL;DR** flow: a **Go** API ([`tldr-ai-be`](tldr-ai-be/README.md)) calls Claude to return a summary and three action items, and a **React Native** app ([`tldr-ai-fe`](tldr-ai-fe/README.md)) pastes text, calls the API, and shows the result.

## Demo

![Screen recording: paste text, tap Summarize, view summary and next steps; usage updates after each run.](demo.gif)

## Environment files

- **`tldr-ai-be/.env.example`** — committed template with safe defaults and every supported variable. No real secrets.
- **Local `.env`** — copy from the example, add `ANTHROPIC_API_KEY` and other secrets. **Never commit** `.env`; it is ignored at the repo root and under `tldr-ai-be/` (see `.gitignore` files).

At startup, the server loads only `.env` files (repo root then `tldr-ai-be/.env`); see `tldr-ai-be/README.md`. `.env.example` is a copy template only.

## Where to go next

| Project | README |
|---------|--------|
| Backend API | [`tldr-ai-be/README.md`](tldr-ai-be/README.md) |
| React Native app | [`tldr-ai-fe/README.md`](tldr-ai-fe/README.md) |
