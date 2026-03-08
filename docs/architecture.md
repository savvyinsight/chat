# Chat Application - Architecture

## System Overview

A real-time chat application built with Go backend and React/TypeScript frontend, enabling direct messaging between users with WebSocket support.

```
┌─────────────────────────────────────────────────────────────────┐
│                     Frontend (React/TypeScript)                  │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │   Login      │  │   Chat       │  │   Me (Profile)
│  │  Component   │  │  Component   │  │  Component   │            │
│  └──────────────┘  └──────────────┘  └──────────────┘            │
│         │                  │                  │                  │
│         └──────────────┬───┴──────────────┬───┘                  │
│                        │                  │                      │
│         ┌──────────────▼──┐    ┌─────────▼──────────┐            │
│         │   REST API      │    │  WebSocket Client  │            │
│         │   (api.ts)      │    │    (ws.ts)         │            │
│         └──────────────┬──┘    └─────────┬──────────┘            │
└────────────────────────┼────────────────┼──────────────────────┘
                         │                │
                    HTTP │                │ WS
                         │                │
┌────────────────────────┼────────────────┼──────────────────────┐
│                        │                │                      │
│            ┌───────────▼────────────────▼──────┐                │
│            │   Gin Web Framework Router        │                │
│            │  (Port 8080)                      │                │
│            └──┬────────────────────────────┬───┘                │
│               │                            │                    │
│    ┌──────────▼───────────────┐   ┌──────▼─────────────┐       │
│    │   REST API Endpoints     │   │  WebSocket Handler │       │
│    │                          │   │                    │       │
│    │ • Auth (Login/Register)  │   │  • ServeWS()       │       │
│    │ • Users (CRUD)           │   │  • Hub Management  │       │
│    │ • Avatar Upload          │   │  • Message Routing │       │
│    │ • Messages (Query)       │   │                    │       │
│    └──────────┬───────────────┘   └────────┬───────────┘       │
│               │                            │                    │
│    ┌──────────▼──────────────────────────────▼────┐             │
│    │        Middleware Layer                      │             │
│    │                                              │             │
│    │ • JWT Authentication                        │             │
│    │ • Token Validation                          │             │
│    │                                              │             │
│    └─────────────┬───────────────────────────────┘             │
│                  │                                              │
│    ┌─────────────▼────────────────────────────┐                │
│    │       Service Layer                      │                │
│    │                                          │                │
│    │ • UserService                            │                │
│    │   - CreateUser, UpdateUser               │                │
│    │   - AuthenticateUser, GetUserList        │                │
│    │                                          │                │
│    │ • MessageService                         │                │
│    │   - SaveMessage, AckMessage              │                │
│    │   - GetMessagesBetween                   │                │
│    │                                          │                │
│    │ • WebSocket Hub                          │                │
│    │   - Manages client connections           │                │
│    │   - Routes messages by user/room         │                │
│    │   - Redis pub/sub coordination           │                │
│    │                                          │                │
│    └─────────────┬────────────────────────────┘                │
│                  │                                              │
│    ┌─────────────▼────────────────────────────┐                │
│    │       Data Layer (GORM ORM)              │                │
│    │                                          │                │
│    │ • Models:                                │                │
│    │   - UserBasic                            │                │
│    │   - Message                              │                │
│    │                                          │                │
│    └─────────────┬────────────────────────────┘                │
│  Backend Go      │                                              │
│  (Port 8080)     │                                              │
└──────────────────┼──────────────────────────────────────────────┘
                   │
        ┌──────────┴──────────┬──────────────┐
        │                     │              │
    MySQL              Redis               File Storage
    (DB)           (Cache/PubSub)        (/static/avatars)
    
   user_basic      • User sessions        • Avatar images
   messages        • Message queues
                   • Broadcast channels
```

## Backend Architecture (Go)

### Project Structure

```
server/
├── api/               # HTTP request handlers
│   ├── index.go      # Welcome endpoint
│   ├── user.go       # User CRUD operations
│   ├── message.go    # Message retrieval
│   └── avatar.go     # Avatar upload
├── service/          # Business logic layer
│   ├── user.go       # User operations
│   └── message.go    # Message operations
├── model/            # GORM data models
│   ├── user_basic.go
│   └── message.go
├── middleware/       # HTTP middleware
│   └── jwt.go       # JWT authentication
├── router/           # Route registration (Gin)
│   └── user_basic.go
├── ws/               # WebSocket handling
│   ├── handler.go   # Connection upgrade
│   ├── client.go    # Client connection
│   ├── hub.go       # Hub management & routing
│   └── hub_test.go
├── config/           # Configuration
│   ├── config.go
│   └── gorm_mysql.go
├── initialize/       # App initialization
│   └── init.go      # Bootstrap routines
├── global/           # Global state
│   └── global.go    # DB/Redis instances
└── main.go          # Entry point
```

### Key Components

#### 1. HTTP API Layer
- **Router**: Gin framework with JWT middleware
- **Endpoints**: RESTful CRUD for users, messages, avatar uploads
- **Authentication**: JWT tokens with 1-hour expiration

#### 2. WebSocket Layer
- **Handler** (`ServeWS`): Upgrades HTTP to WebSocket, validates auth
- **Client**: Maintains connection, read/write pumps, message buffering
- **Hub**: Routes messages, manages connections per user/room
- **Pub/Sub**: Redis integration for distributed messaging

#### 3. Service Layer
- Encapsulates business logic
- User authentication, validation, CRUD operations
- Message persistence and retrieval

#### 4. Data Layer
- GORM ORM for database abstraction
- MySQL connection pooling
- Automatic timestamps (created_at, updated_at)

### Data Flow

#### Direct Message Flow (WebSocket)
```
Frontend App
    │
    ├─► Establish WS connection (/ws?token=...) ► ServeWS()
    │                                                │
    │                                            Upgrade HTTP
    │                                                │
    │                                            Create Client
    │                                                │
    │                                            Register in Hub
    │
    ├─► Send message (JSON) ► Client.ReadPump() ► Message processing
    │                                                │
    │                                            Hub.Broadcast
    │                                                │
    │                                            Route by user/room
    │                                                │
    │                                            SaveMessage() ► DB
    │
    └─ Receive message ◄ Client.WritePump() ◄ Hub routes to clients
```

#### REST API Flow (Login & Message History)
```
Frontend
    │
    ├─► POST /user/login (email/phone + password)
    │       │
    │       └─► api.Login() ► service.AuthenticateUser()
    │                            │
    │                            ├─► Verify credentials
    │                            └─► Generate JWT token
    │
    ├─► GET /user/me (with JWT)
    │       │
    │       └─► api.GetCurrentUser() ► Verify JWT middle ware
    │                                      │
    │                                      continue...