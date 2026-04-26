# Tldr-ai frontend (React Native CLI)

**TldrAi** — minimal UI: paste text, call the Go backend, show a **summary** and **three action items**.

This project uses **Yarn** (see `packageManager` in [`package.json`](package.json)). Install uses a classic **`node_modules`** layout ([`.yarnrc.yml`](.yarnrc.yml) `nodeLinker: node-modules`) so Metro and React Native work reliably.

## Requirements

- [Corepack](https://nodejs.org/api/corepack.html) enabled (ships with Node 16.10+): `corepack enable` — then Yarn version follows `packageManager` in `package.json`.
- [React Native dev environment](https://reactnative.dev/docs/set-up-your-environment) (Metro, Xcode and/or Android toolchain).
- **Node** — **≥ 22.11.0 required** (`package.json` `engines`). Metro depends on **`Array.prototype.toReversed()`** (ES2023); older Node (e.g. 18) fails with `configs.toReversed is not a function` when you run `yarn start`. Use [`.nvmrc`](.nvmrc) / [`.node-version`](.node-version) with **nvm** / **fnm** (`nvm install && nvm use`), or **asdf**: `asdf install nodejs 22.11.0` (see [`.tool-versions`](.tool-versions) at the app root — `ios/.tool-versions` is Ruby-only for CocoaPods).
- **Backend** — run [`tldr-ai-be`](../tldr-ai-be/README.md) first; the app expects the API on port **8080** by default.

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

**iOS (CocoaPods):** from `ios/`:

```bash
cd ios
bundle install
bundle exec pod install
```

If you use **[asdf](https://asdf-vm.com/)**, Ruby is pinned in [`.tool-versions`](.tool-versions) and [`ios/.tool-versions`](ios/.tool-versions) (default **3.2.3**). Install it if needed: `asdf install ruby 3.2.3`, or change both files to another version you already have (e.g. `3.4.4`) so `bundle` is not blocked by “No version is set for command bundle”.

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
