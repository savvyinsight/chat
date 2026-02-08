import React, { useEffect, useState } from 'react'
import Login from './components/Login'
import Chat from './components/Chat'
import Me from './components/Me'
import api from './api'

type User = any

export default function App() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))
  const [userId, setUserId] = useState<string | null>(() => localStorage.getItem('user_id'))
  const [tab, setTab] = useState<'friends' | 'group' | 'me'>('friends')
  const [me, setMe] = useState<User | null>(null)

  const onLogin = (t?: string, uid?: string) => {
    if (t) {
      localStorage.setItem('token', t)
      setToken(t)
    }
    if (uid) {
      localStorage.setItem('user_id', uid)
      setUserId(uid)
    }
  }

  const onLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user_id')
    setToken(null)
    setUserId(null)
  }

  useEffect(() => {
    if (!token) return
    api.get('/user/me').then((b) => setMe(b.data)).catch(() => setMe(null))
  }, [token])

  if (!token && !userId) {
    return <Login onLogin={onLogin} />
  }

  return (
    <div className="h-screen flex flex-col">
      <header className="bg-white p-3 shadow flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="text-xl font-bold">Chat</div>
          {me && (
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <img src={me.AvatarURL || me.avatar_url || '/static/avatars/default.png'} className="w-8 h-8 rounded-full object-cover" />
              <div>{me.Name || me.name || 'Me'}</div>
            </div>
          )}
        </div>
        <div className="flex items-center gap-3">
          <nav className="space-x-3">
            <button className={`px-3 py-1 ${tab === 'friends' ? 'bg-gray-200' : ''}`} onClick={() => setTab('friends')}>Friends</button>
            <button className={`px-3 py-1 ${tab === 'group' ? 'bg-gray-200' : ''}`} onClick={() => setTab('group')}>Group</button>
            <button className={`px-3 py-1 ${tab === 'me' ? 'bg-gray-200' : ''}`} onClick={() => setTab('me')}>Me</button>
          </nav>
          <button className="text-sm text-red-500" onClick={onLogout}>Logout</button>
        </div>
      </header>
      <main className="flex-1 overflow-auto bg-gray-100">
        {tab === 'friends' && <Chat token={token} userId={userId ? Number(userId) : undefined} onLogout={onLogout} />}
        {tab === 'group' && <div className="p-6">Group chat (coming soon)</div>}
        {tab === 'me' && <Me onProfileUpdate={(u) => setMe(u)} />}
      </main>
    </div>
  )
}
