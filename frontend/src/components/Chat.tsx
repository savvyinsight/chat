import React, { useEffect, useRef, useState } from 'react'
import { createSocket, Message } from '../ws'
import api from '../api'

export default function Chat({ token, userId, onLogout }: { token?: string | null; userId?: number; onLogout: () => void }) {
  const [users, setUsers] = useState<Array<any>>([])
  const [selectedUser, setSelectedUser] = useState<number | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    api.get('/userList').then((b) => setUsers(b.data || [])).catch(() => setUsers([]))
  }, [])

  useEffect(() => {
    const ws = createSocket({ token: token || undefined, userId: userId })
    wsRef.current = ws
    ws.onopen = () => console.log('ws open')
    ws.onmessage = (ev) => {
      try {
        const m: Message = JSON.parse(ev.data)
        setMessages((prev) => [...prev, m])
      } catch (e) {
        console.warn('invalid message', e)
      }
    }
    ws.onclose = () => console.log('ws closed')
    ws.onerror = (e) => console.error('ws error', e)
    return () => ws.close()
  }, [token, userId])

  const send = () => {
    if (!wsRef.current) return
    if (!selectedUser) return
    const m: Message = { type: 'message', to: selectedUser, body: input }
    wsRef.current.send(JSON.stringify(m))
    setMessages((prev) => [...prev, { ...m, from: userId }])
    setInput('')
  }

  return (
    <div className="p-6 h-screen grid grid-cols-4 gap-4">
      <div className="col-span-1 bg-white p-4 rounded shadow">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold">Users</h3>
          <button className="text-red-500 text-sm" onClick={onLogout}>Logout</button>
        </div>
        <div className="space-y-2">
          {users.map((u: any) => (
            <div key={u.ID || u.id} className={`p-2 rounded cursor-pointer ${selectedUser === u.ID ? 'bg-blue-100' : 'hover:bg-gray-100'}`} onClick={() => setSelectedUser(u.ID)}>
              <div className="font-semibold">{u.Name || u.name || `user:${u.ID}`}</div>
              <div className="text-xs text-gray-500">id: {u.ID}</div>
            </div>
          ))}
        </div>
      </div>
      <div className="col-span-3 bg-white p-4 rounded shadow flex flex-col">
        <div className="flex-1 overflow-auto mb-4">
          {messages.map((m, i) => (
            <div key={i} className={`mb-2 p-2 rounded ${m.from === userId ? 'bg-blue-50 self-end' : 'bg-gray-100 self-start'}`}>
              <div className="text-sm">{m.body}</div>
              <div className="text-xs text-gray-500">{m.type} {m.id ? `#${m.id}` : ''} {m.from ? `from ${m.from}` : ''}</div>
            </div>
          ))}
        </div>
        <div className="flex gap-2">
          <input className="flex-1 p-2 border rounded" value={input} onChange={(e) => setInput(e.target.value)} placeholder={selectedUser ? 'Type a message' : 'Select a user to message'} />
          <button className="px-4 py-2 bg-green-600 text-white rounded" onClick={send} disabled={!selectedUser}>Send</button>
        </div>
      </div>
    </div>
  )
}
