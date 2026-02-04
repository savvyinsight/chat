# Chat Frontend

This is a minimal Vite + React + TypeScript + Tailwind frontend for the chat server.

Quick start:

1. cd frontend
2. npm install
3. npm run dev

Features included:
- Login form (posts to `/user/login`; middleware JWT `token` is preferred)
- User list fetched from `/userList`
- WebSocket connection to `/ws` (passes `token` query param when available, or `user_id` fallback)

Next steps:
- Add better message history API calls and UI
- File upload support via REST
- Serve built `dist` from Go server (let me know if you want that wired up)
