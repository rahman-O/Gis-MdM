# Gis-MdM — Mobile Device Management Platform

نظام إدارة أجهزة محمولة (MDM) متكامل مبني على Headwind MDM مع واجهة React حديثة وباكند Go عالي الأداء.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Clients                               │
├──────────────────┬──────────────────┬───────────────────────┤
│  React Frontend  │  Flutter Agent   │   Legacy HMDM Agent   │
│  (Admin Panel)   │  (Device Side)   │   (Android APK)       │
│  Port: 5173      │  Background Svc  │                       │
└────────┬─────────┴────────┬─────────┴───────────┬───────────┘
         │                  │                     │
         ▼                  ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Backend (API)                           │
│                    Port: 8081                                 │
│                                                              │
│  /rest/private/*  — Authenticated admin endpoints            │
│  /rest/public/*   — Device sync & enrollment endpoints       │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  PostgreSQL 14                                │
│                  Port: 5432                                   │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
Gis-MdM/
├── frontend/              # React + Vite + TypeScript (Admin UI)
├── serverBackendGo/       # Go backend (Gin framework)
├── flutter_mdm_agent/     # Flutter MDM Agent (Android)
├── docker-compose.yml     # Production deployment
└── Makefile               # Root-level commands
```

## Quick Start (Development)

### Prerequisites
- Go 1.22+
- Node.js 20+
- Flutter 3.11+
- Docker & Docker Compose
- PostgreSQL 14 (via Docker)

### 1. Start Database
```bash
cd serverBackendGo
docker compose up -d   # Starts PostgreSQL on port 5432
```

### 2. Start Backend
```bash
cd serverBackendGo
make dev               # Runs migrations + starts Go server on :8081
```

### 3. Start Frontend
```bash
cd frontend
npm install
npm run dev            # Vite dev server on :5173 (proxies /rest → :8081)
```

### 4. Build Flutter Agent
```bash
cd flutter_mdm_agent
flutter pub get
flutter build apk --debug
```

### Install on Device
```bash
adb install -r flutter_mdm_agent/build/app/outputs/flutter-apk/app-debug.apk
```

## Key Features

- **Device Management** — Enroll, monitor, and control Android devices
- **Real-time Telemetry** — Battery, location, network, storage, system info
- **Background Service** — Agent runs permanently (Foreground Service + Device Owner)
- **QR Enrollment** — Factory reset → scan QR → fully managed device
- **Policy Engine** — Kiosk mode, app restrictions, hardware controls
- **Device Tree** — Hierarchical folder organization for devices
- **Profile System** — Configuration profiles with versioning and rollout
- **Enrollment Routes** — Multiple enrollment paths with different configs

## Environment Variables

### Backend (`serverBackendGo/.env`)
| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://hmdm:hmdm@localhost:5432/hmdm?sslmode=disable` |
| `SERVER_PORT` | HTTP server port | `8081` |
| `JWT_SECRET` | JWT signing key | (required) |
| `HASH_SECRET` | Sync response signature | (required) |

## Deployment

### Docker (Production)
```bash
docker compose up -d --build
```
- Frontend: http://localhost:3000
- Backend API: http://localhost:8081
- Database: localhost:5433

### HTTPS (studhub.app)
```bash
make setup-studhub-https   # Configure Cloudflare tunnel
make dev-https             # Start with HTTPS
```

## Flutter Agent — Device Enrollment

### QR Code Provisioning (Recommended)
1. Factory reset the device
2. On welcome screen, tap 6 times on empty area
3. Scan the enrollment QR code
4. Device auto-provisions with MDM Agent as Device Owner

### Manual Installation (Development)
```bash
adb install -r flutter_mdm_agent/build/app/outputs/flutter-apk/app-debug.apk
adb shell am start -n com.gismdm.mdm_agent/.MainActivity
```

## API Documentation

Swagger UI available at: `http://localhost:8081/swagger/index.html`

## License

Private — All rights reserved.
