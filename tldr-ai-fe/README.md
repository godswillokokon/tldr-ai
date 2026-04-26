# Tldr-ai frontend (React Native CLI)

**TldrAi** — paste or type text, call the Go backend, show a **summary** and **three action items**.

This project uses **Yarn** (see `packageManager` in [`package.json`](package.json)). Install uses a classic **`node_modules`** layout ([`.yarnrc.yml`](.yarnrc.yml) `nodeLinker: node-modules`) so Metro and React Native work reliably.

## Prerequisites

- **[Corepack](https://nodejs.org/api/corepack.html)** — `corepack enable`; Yarn version follows `packageManager` in `package.json`.
- **[React Native dev environment](https://reactnative.dev/docs/set-up-your-environment)** — Metro, Xcode (iOS), and/or Android Studio / SDK.
- **Node ≥ 22.11.0** — `package.json` `engines`; use [`.nvmrc`](.nvmrc) / [`.node-version`](.node-version) (nvm, fnm, or asdf for Node).
- **Backend on port 8080** — start [`tldr-ai-be`](../tldr-ai-be/README.md) before exercising the app from a simulator or emulator.
- **iOS (CocoaPods)** — from `ios/`: `bundle install` then `bundle exec pod install`. Ruby is pinned in [`.tool-versions`](.tool-versions) and [`ios/.tool-versions`](ios/.tool-versions) for asdf users.

## Quick smoke test (curl)

With **tldr-ai-be** listening on `8080` (default `PORT`):

```bash
# Health
curl -sS "http://localhost:8080/health"

# Usage snapshot (JSON)
curl -sS "http://localhost:8080/api/usage"

# Process text (needs ≥20 runes of text on the server; example uses 20+ ASCII chars)
curl -sS -X POST "http://localhost:8080/api/processText" \
  -H "Content-Type: application/json" \
  -d '{"text":"aaaaaaaaaaaaaaaaaaaa"}'
```

Expect `200` and JSON with `summary` and `actionItems` (three strings) when the AI provider is configured. A `503` or error JSON usually means missing/invalid `ANTHROPIC_API_KEY` on the server.

## Frontend ↔ backend alignment

- **Minimum input** — the app enables submit when **trimmed length ≥ 20** characters ([`MIN_INPUT_LENGTH`](src/domain/constants.ts)). The Go API enforces **≥ 20 UTF-8 runes** after trim; for most English text these match. Very high Unicode-to-byte edge cases can differ slightly.
- **Success JSON** — `POST /api/processText` returns **`summary`** (string) and **`actionItems`** (array of exactly three strings), optional **`model`**. The client parser rejects anything else ([`parseProcessTextResponse`](src/domain/parseProcessResponse.ts)).
- **Usage** — `GET /api/usage` returns the Go `Snapshot` shape (`used`, `cap`, `unlimited`, `spentUsd`, `budgetUsd`, …). The app maps that in [`parseUsageResponse`](src/domain/parseUsageResponse.ts); set caps with **`USAGE_MAX_CALLS`** / **`USAGE_BUDGET_USD`** in the server `.env` (see `tldr-ai-be/.env.example`).
- **Android HTTP** — debug builds set **`usesCleartextTraffic`** to **`true`** via Gradle `manifestPlaceholders` so `http://10.0.2.2:8080` reaches your host. **Release** uses **`false`**; ship **HTTPS** or a network security config if you need cleartext in production.

## Install and run

```bash
yarn install
yarn start
```

In another terminal:

```bash
yarn ios
# or
yarn android
```

```bash
yarn lint
yarn test
```

`yarn test` runs Jest suites including `parseProcessResponse` (pure parser), `fetchProcessText` (fetch mock), and the smoke `App` render.

## How it talks to the backend

[`src/config/apiConfig.ts`](src/config/apiConfig.ts) builds the process URL:

| Runtime | Base URL | Notes |
|---------|-----------|--------|
| **iOS Simulator** | `http://localhost:8080` | Host Mac loopback |
| **Android Emulator** | `http://10.0.2.2:8080` | Maps to host machine (not `localhost` on the emulator) |

Path: **`/api/processText`** (see `PROCESS_TEXT_PATH` and `getProcessTextUrl()`).

**Physical device:** point the base URL at your computer’s LAN IP (e.g. `http://192.168.1.x:8080`), or use something like `adb reverse tcp:8080 tcp:8080` on Android and keep a host-aligned URL—adjust `getApiBaseUrl()` / constants accordingly.

## Source layout

| Path | Role |
|------|------|
| [`src/domain/types.ts`](src/domain/types.ts) | `ProcessResponse` type |
| [`src/domain/constants.ts`](src/domain/constants.ts) | `MIN_INPUT_LENGTH` (20; keep in sync with backend) |
| [`src/domain/parseProcessResponse.ts`](src/domain/parseProcessResponse.ts) | Parses `text` response safely (JSON + shape checks) |
| [`src/config/apiConfig.ts`](src/config/apiConfig.ts) | Port, path, `getApiBaseUrl()`, `getProcessTextUrl()` |
| [`src/data/processTextRepository.ts`](src/data/processTextRepository.ts) | `fetch` + `parseProcessTextResponse`; maps network failures to clear errors |
| [`src/hooks/useProcessText.ts`](src/hooks/useProcessText.ts) | Text state, loading, submit, alerts; optional `processFn` for tests |
| [`src/components/`](src/components/) | `TldrHeader`, `TextBlockInput`, `ProcessTextButton`, `ProcessResultCard` |
| [`App.tsx`](App.tsx) | `SafeAreaProvider` + composes hook + components (uses `react-native-safe-area-context`) |

## Behaviour

- Submit is enabled when trimmed input length is **≥ `MIN_INPUT_LENGTH`** and not loading.
- Errors from the API (`error` field in JSON) are shown in an alert.
- Successful responses render summary, numbered action items, and `model` when present.
