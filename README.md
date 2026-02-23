# Chat — Self-Hosted Real-Time Messaging Backend

A production-ready, self-hosted chat backend built with Go (Gin), providing REST APIs and WebSocket real-time messaging. Designed for teams needing private, scalable in-app messaging or chat features.

## Features
- **Real-Time Messaging**: WebSocket hub (`/ws`) with per-user and group chat routing, heartbeat/ping for reliability, and Redis pub/sub for horizontal scaling.
- **Authentication**: JWT-based auth with endpoints for login (`POST /user/login`) and registration (`POST /user/register`).
- **User Management**: Profile retrieval (`GET /user/me`), avatar uploads (`POST /user/avatar`), and static file serving.
- **Message Persistence**: History API (`GET /messages`), GORM + MySQL storage, and message receipts (planned).
- **Scalability**: In-memory hub with Redis for cross-instance pub/sub; supports group rooms and file attachments (planned).

## Architecture
- **Server**: Gin router with middleware for JWT auth, CORS, and logging.
- **WebSocket Hub**: Manages client connections, routes messages (1-on-1 or group), and handles reconnections.
- **Persistence**: GORM models for users and messages; Redis for pub/sub.
- **Config**: YAML-based configuration for DB, Redis, and server settings.

## Quick Start (Backend Focus)
1. **Prerequisites**: Go 1.21+, MySQL, Redis (optional for scaling).
2. **Clone & Setup**:
   ```bash
   git clone <repo>
   cd chat/server
   go mod tidy
   cp config.example.yaml config.yaml  # Edit DB/Redis settings