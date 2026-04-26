# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go, module: `Diggpher`, go 1.24.0)
- **Start dev server**: `cd backend/cmd && go run main.go` (listens on `:9090`)
- **Add dependencies**: `cd backend && go get <pkg>` then `go mod tidy`
- **Build binary**: `cd backend/cmd && go build -ldflags "-s -w" -o watchauth-backend.exe`

### Frontend (Vue 3 + Vite, listens on `:3007`)
- **Start dev server**: `cd frontend && npm run dev`
- **Install deps**: `npm install` (in frontend/)
- **Production build**: `cd frontend && npm run build-only`
- **Type check**: `cd frontend && npm run type-check`
- **Lint**: `cd frontend && npm run lint`
- **Format**: `cd frontend && npm run format`

### Full Build
- `python build.py` from root — builds Go binary + frontend, copies configs to `./dist/`

## Architecture

### Backend Layering (`backend/`)
```
cmd/main.go → initialize/* → internal/route/ → internal/controller/ → internal/service/ → internal/dao/
```
- **`initialize/`** — Bootstrap: config loading (viper from `configs/config.yaml`), PostgreSQL (GORM), Redis, Fiber app with CORS, then route binding
- **`global/`** — Package-level singletons: `CONFIG`, `DataBase` (*gorm.DB), `Redis` (*redis.Client), `WebApp` (*fiber.App), `Log`/`SugarLog` (zap)
- **`internal/controller/`** — Fiber handlers, parse request, call service, respond via `Respond()` helper
- **`internal/service/`** — Business logic, coordinates DAOs and Redis cache
- **`internal/dao/`** — GORM model definitions + AutoMigrate in `BindDao()`
- **`internal/service/errMsg/`** — Centralized error codes and Chinese error messages
- **`pkg/middleware/auth/`** — JWT generation and Bearer token middleware
- **`pkg/logger/`** — Zap + lumberjack (console + file rotation)

### Frontend (`frontend/src/`)
- **`api/`** — Axios-based API modules (one file per domain: login, app, version, user, role, menu)
- **`stores/`** — Pinia stores: `menu` (sidebar state + icons), `user` (profile + messages), `tabs` (tab navigation), `theme` (dark/light + primary color)
- **`router/`** — Vue Router with dynamic route generation from menu permissions
- **`types/`** — TypeScript interfaces for all API responses
- Uses **MSW** (Mock Service Worker) for development — enable/disable via `APP_CONFIG.enableMSW`

### Key Routes
- `/api/admin/login` (POST, no auth) — Admin login
- `/api/admin/*` (with JWT middleware) — App CRUD, version management, machine queries
- `/api/sdk/login` (POST, no auth) — Machine SDK authentication

### Database Models (`internal/dao/`)
- `Admin` — Admin users
- `App` — Managed software applications
- `Version` — App version releases
- `Machine`, `MachineInfo`, `UsedApp`, `MachineLog` — Machine registration, metadata, app usage tracking

### Configuration
- Backend: `backend/configs/config.yaml` (DB, Redis, server port)
- Backend config struct in `global/config.go` — also reads Logger config (level, console, dir)
- Frontend: `.env.development` (`VITE_API_BASE_URL`, `VITE_STATIC_URL`) and `.env.production`
- Frontend app config: `src/config/app.config.ts` (MSW toggle, project name, theme toggle)
