import React, { useState } from 'react'
import Login from './components/Login'
import Chat from './components/Chat'

export default function App() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))
  const [userId, setUserId] = useState<string | null>(() => localStorage.getItem('user_id'))

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

  if (!token && !userId) {
    return <Login onLogin={onLogin} />
  }

  return <Chat token={token} userId={userId ? Number(userId) : undefined} onLogout={onLogout} />
}
