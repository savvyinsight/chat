# Chat Application - Database Schema

## Overview

The application uses MySQL with GORM ORM. Database is configured via `config.yaml` with automatic migration on startup.

```yaml
DNS: root:root@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local
```

## Data Models

### 1. User Model (`user_basic` table)

#### Purpose
Stores user authentication and profile information.

#### Schema

| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| `id` | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | User identifier |
| `created_at` | TIMESTAMP | - | Record creation time (GORM auto) |
| `updated_at` | TIMESTAMP | - | Last modification time (GORM auto) |
| `deleted_at` | TIMESTAMP | NULL | Soft delete (GORM auto) |
| `name` | VARCHAR(255) | - | Display name |
| `avatar_url` | VARCHAR(500) | - | URL to user avatar image |
| `password` | VARCHAR(255) | - | Bcrypt hashed password |
| `phone` | VARCHAR(20) | UNIQUE INDEX | Phone number (used for login/registration) |
| `email` | VARCHAR(255) | UNIQUE INDEX | Email address (used for login/registration) |
| `identity` | VARCHAR(255) | - | Additional identity info (reserved) |
| `client_ip` | VARCHAR(45) | - | Last login IP address |
| `client_port` | VARCHAR(10) | - | Last login port |
| `login_time` | BIGINT UNSIGNED | - | Last login timestamp (Unix) |
| `heartbeat_time` | BIGINT UNSIGNED | - | Last activity timestamp |
| `logout_time` | BIGINT UNSIGNED | - | Last logout timestamp |
| `is_logout` | BOOLEAN | DEFAULT FALSE | Current logout status |
| `device_info` | VARCHAR(500) | - | Last login device metadata |

#### Go Model Definition
```go
type UserBasic struct {
    gorm.Model
    Name          string
    AvatarURL     string
    PassWord      string
    Phone         string `gorm:"uniqueIndex"`
    Email         string `gorm:"uniqueIndex"`
    Identity      string
    ClientIp      string
    ClientPort    string
    LoginTime     uint64
    HeartbeatTime uint64
    LogoutTime    uint64
    IsLogout      bool
    DeviceInfo    string
}
```

#### Indexes
- `idx_user_basic_phone` (UNIQUE) - Fast lookup by phone during login
- `idx_user_basic_email` (UNIQUE) - Fast lookup by email during login
- `idx_user_basic_deleted_at` - GORM soft delete filtering

#### Constraints
- **Email**: Must be unique or empty; validated with govalidator
- **Phone**: Must be unique or empty; validated with govalidator
- **Soft Delete**: Uses `deleted_at` for safe data retention
- **At least one identifier required**: Email OR phone for authentication

---

### 2. Message Model (`messages` table)

#### Purpose
Stores chat messages with delivery metadata for direct and room-based messaging.

#### Schema

| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| `id` | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | Message identifier |
| `created_at` | TIMESTAMP | - | Message creation time |
| `updated_at` | TIMESTAMP | - | Last modification time |
| `deleted_at` | TIMESTAMP | NULL | Soft delete |
| `from` | BIGINT UNSIGNED | INDEX | Sender user ID (FK: user_basic.id) |
| `to` | BIGINT UNSIGNED | INDEX | Recipient user ID (FK: user_basic.id) |
| `room` | VARCHAR(255) | INDEX | Room identifier (for group chat) |
| `type` | VARCHAR(50) | - | Message type (e.g., "text", "image", "voice") |
| `body` | LONGTEXT | - | Message content (supports JSON for rich content) |
| `delivered` | BOOLEAN | DEFAULT FALSE | Delivery status flag |
| `delivered_at` | TIMESTAMP | NULL | Delivery confirmation time |

#### Go Model Definition
```go
type Message struct {
    gorm.Model
    From        uint       `json:"from" gorm:"index"`
    To          uint       `json:"to,omitempty" gorm:"index"`
    Room        string     `json:"room,omitempty" gorm:"index"`
    Type        string     `json:"type"`
    Body        string     `json:"body" gorm:"type:text"`
    Delivered   bool       `json:"delivered"`
    DeliveredAt *time.Time `json:"delivered_at"`
}
```

#### Indexes
- `idx_messages_from` - Filter messages sent by a user
- `idx_messages_to` - Filter messages received by a user
- `idx_messages_room` - Filter messages in a room
- `idx_messages_deleted_at` - GORM soft delete filtering

#### Constraints
- **Foreign Key** (Recommended): `from` → `user_basic.id`
- **Foreign Key** (Recommended): `to` → `user_basic.id`
- **Message Type**: Use constants ("text", "image", "video", etc.)
- **At least one recipient required**: Either `to` (direct) or `room` (group)

#### Usage Patterns
- **Direct message**: `from` + `to` populated, `room` NULL
- **Group message**: `from` + `room` populated, `to` NULL
- **Message history query**: 
  ```sql
  WHERE (from = ? AND to = ?) OR (from = ? AND to = ?)
  ORDER BY id ASC
  LIMIT 100
  ```

---

## Data Relationships

```
┌───────────────────┐
│    user_basic     │
│                   │
│ id (PRIMARY KEY)  │
│ email (UNIQUE)    │
│ phone (UNIQUE)    │
│ name              │
│ ...               │
└─────────┬─────────┘
          │
          │ 1:N (sender)
          ├─────────────────┐
          │                 │
          │ 1:N (recipient) │
          │                 │
      ┌───▼──────────────┐
      │   messages       │
      │                  │
      │ id (PK)          │
      │ from (FK, INDEX) │──┐
      │ to (FK, INDEX)   │──┤ references user_basic(id)
      │ room (INDEX)     │
      │ body             │
      │ delivered        │
      │ ...              │
      └──────────────────┘
```

---

## Database Operations

### User Operations

#### Create User
```sql
INSERT INTO user_basic (created_at, updated_at, name, email, phone, password)
VALUES (NOW(), NOW(), ?, ?, ?, ?);
```

#### Authenticate User
```sql
SELECT * FROM user_basic 
WHERE (email = ? OR phone = ?) AND deleted_at IS NULL;
```

#### Update User Profile
```sql
UPDATE user_basic 
SET name = ?, avatar_url = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;
```

#### Get User by ID
```sql
SELECT * FROM user_basic 
WHERE id = ? AND deleted_at IS NULL;
```

#### Delete User (Soft Delete)
```sql
UPDATE user_basic 
SET deleted_at = NOW()
WHERE id = ?;
```

### Message Operations

#### Save Message
```sql
INSERT INTO messages (created_at, updated_at, from, to, room, type, body, delivered)
VALUES (NOW(), NOW(), ?, ?, ?, ?, ?, FALSE);
```

#### Mark Message as Delivered
```sql
UPDATE messages 
SET delivered = TRUE, delivered_at = NOW()
WHERE id = ?;
```

#### Get Message History (Direct)
```sql
SELECT * FROM messages 
WHERE (from = ? AND to = ?) OR (from = ? AND to = ?)
AND deleted_at IS NULL
ORDER BY id ASC
LIMIT 100;
```

#### Get Room Messages
```sql
SELECT * FROM messages 
WHERE room = ? AND deleted_at IS NULL
ORDER BY id ASC
LIMIT 100;
```

---

## Indexes & Performance

### Query Performance Strategy

| Query | Index | Lookup Time |
|-------|-------|------------|
| Find user by email | email (UNIQUE) | O(1) |
| Find user by phone | phone (UNIQUE) | O(1) |
| Find messages for user | from, to (INDEX) | O(log N) |
| Find messages in room | room (INDEX) | O(log N) |
| List all users | None (full scan) | O(N) |

### Recommended Additional Indexes (Future)

```sql
-- Composite index for message history queries
CREATE INDEX idx_messages_from_to ON messages(from, to, created_at DESC);
CREATE INDEX idx_messages_room_created ON messages(room, created_at DESC);

-- Index for delivered status (useful for delivery confirmation queries)
CREATE INDEX idx_messages_undelivered ON messages(delivered, id);

-- Index for user activity queries
CREATE INDEX idx_user_heartbeat ON user_basic(heartbeat_time DESC);
```

---

## Configuration

### MySQL Connection Parameters (config.yaml)

```yaml
Mysql:
  dns: "root:root@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"
```\n\n#### DSN Breakdown
- `charset=utf8mb4` - Full Unicode support (emoji, Chinese)
- `parseTime=True` - Auto-convert TIMESTAMP to Go time.Time
- `loc=Local` - Use server timezone

### Initialization
- Database auto-creates on first connection
- GORM auto-migration runs on startup (`config/gorm_mysql.go`)
- Connection pooling defaults: max 10 open connections

---

## Future Schema Enhancements

| Feature | Tables/Columns | Status |
|---------|----------------|--------|
| User blocking | Add `user_blocks` table | Planned |
| Read receipts | Add `read_at` column to messages | Planned |
| Group chats | Extend `room` usage + `room_members` table | Planned |
| Message reactions | Add `message_reactions` table | Planned |
| User presence | Add `status` column to user_basic | Planned |
| Message encryption | Add `is_encrypted`, `salt` columns | Research needed |
