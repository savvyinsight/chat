# Chat — Self-hosted Real-time Messaging

A lightweight, production-minded example of a self-hosted chat service:

- Backend: Go (Gin) providing REST APIs and a WebSocket real-time hub.
- Frontend: React + TypeScript + Vite + Tailwind for fast development and a modern UI.

This repository is intended as a starting point for teams building a private chat or in-app messaging feature.

Key features
 - WebSocket-based real-time messaging at `/ws` with per-user routing and Redis pub/sub support for horizontal scaling.
 - REST APIs for authentication (`/user/login`), user management, and message history (`/messages`).
 - Developer workflow with Vite dev server proxying API and WS calls for same-origin behavior in development.

Architecture overview
- The server maintains a Hub that registers `Client` connections and routes messages (in-memory) plus publishes messages to Redis so other instances receive them.
- Messages are persisted to a SQL database via GORM; delivery acknowledgements are supported.
- The frontend connects via `/ws` using either a JWT token or a `user_id` query parameter (dev fallback). REST calls use the same JWT for authentication.

Prerequisites
- Go 1.20+ (for the backend)
- Node.js 18+ and npm (for the frontend dev flow)

Quickstart — Development

1) Start the backend (default port 8080):

```bash
cd server
go run main.go
```

2) Start the frontend dev server (Vite, default port 5173):

```bash
cd frontend
npm install
npm run dev
```

Open: http://localhost:5173

Notes
- The Vite dev server proxies API and WebSocket requests to the backend (see `frontend/vite.config.ts`). This avoids CORS during development and lets the frontend call relative paths like `/user/login` and `/ws`.

Authentication
- Login: `POST /user/login` { identifier, password } → returns `{ user_id, token }`.
- The frontend stores the received `token` (JWT) in `localStorage` and attaches it to REST requests (`Authorization: Bearer ...`) and to the WebSocket connection as `?token=...`.

Message history
- Endpoint: `GET /messages?with=<other_user_id>` (authenticated). The server derives the current user from the JWT.

Frontend dev
- Code lives in `frontend/`.
- Key files:
	- `src/components/Login.tsx` — login UI and token handling
	- `src/components/Chat.tsx` — chat UI, history loading, optimistic sends
	- `src/ws.ts` — WebSocket helper with reconnection and buffering

Production build and serving
- Build the frontend:

```bash
cd frontend
npm run build
# outputs a `dist/` directory
```

- Options to serve build artifacts:
	- Serve `dist/` from a static file server (NGINX, S3 + CloudFront, etc.).
	- Embed `dist/` into the Go binary (using `embed`) and serve from `/` so API and UI are same-origin — I can add this integration if you want.

Next recommended improvements
- Harden auth flows and validate tokens on all protected endpoints.
- Add server-driven read/delivery receipts and update optimistic messages accordingly.
- Improve WS reliability: heartbeat/ping, stronger reconnection policies, and message ordering guarantees.
- Add file attachments (upload via REST + message metadata via WS).

Contributing
- Fork, make changes on a branch, run tests and linting, then open a PR.

License
- MIT (or change to your preferred license).

---
For a tailored README with more architecture diagrams or deployment manifests, tell me what target audience (internal devs, OSS users, or ops) you want it oriented to.
