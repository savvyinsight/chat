# Project Progress

This document tracks what has been implemented in the chat project and what remains to be done. It complements `requirements.md` by turning requirements into a checklist and noting current status.

## Completed Features

- ✅ **User registration** (phone/email + password) with JWT token return
- ✅ **User login** with JWT token, token stored and used by frontend
- ✅ **JWT authentication middleware** protecting REST routes
- ✅ **User profile endpoints**: `GET /user/me`, update/patch/delete
- ✅ **Avatar upload and serving** via `/user/avatar` and `/static/avatars`
- ✅ **Message persistence** in MySQL via GORM
- ✅ **Message history API** `GET /messages?with=<id>`
- ✅ **WebSocket hub** with per-user routing, heartbeat/ping, reconnection
- ✅ **Group chat support** (room messages, join/leave, Redis pub/sub)
- ✅ **Swagger/OpenAPI documentation** available at `/swagger`
- ✅ **Frontend components** for login/register/chat/profile
- ✅ **WebSocket client wrapper** with exponential backoff
- ✅ **Chat UI features**: status indicator, message deduplication, history load
- ✅ **Backend tests** covering auth, messaging, WS hub

## In Progress / Partial

- 🟡 **Read/delivery receipts**: database fields added but API not exposed
- 🟡 **Default avatar placeholder**: frontend needs UI & backend default
- 🟡 **Client-side phone validation and redirect**: partially implemented

## Pending Features

- ⬜ Forgot password flow (email or SMS reset)
- ⬜ Message reactions/edit/delete
- ⬜ Media/file attachments via REST + WS metadata
- ⬜ User presence/status indicators
- ⬜ Typing indicator
- ⬜ Message encryption/end-to-end
- ⬜ Contact/friend list and blocking
- ⬜ Read receipts (complete implementation)
- ⬜ Tests for new features (e.g., attachments, receipts)

## Non‑functional & DevOps Tasks

- ⬜ Containerization (Dockerfile, docker-compose)
- ⬜ CI/CD pipeline (GitHub Actions)
- ⬜ Load testing and performance tuning
- ⬜ Monitoring/log aggregation configuration
- ⬜ Deployment documentation and scripts
- ⬜ Secure configuration management

## How to Use This Document

Update this file whenever you finish or start a feature. Each checklist item corresponds to a requirement; move it between sections as you progress. This keeps the team aware of current status and next steps.