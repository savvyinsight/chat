import React, { useState } from 'react'
import api from '../api'
import Register from './Register'

export default function Login({ onLogin }: { onLogin: (token?: string, userId?: string) => void }) {
  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [showRegister, setShowRegister] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    try {
      const body = await api.post('/user/login', { identifier, password })
      // Accept either token (JWT) or user_id fallback
      const token = (body as any).token || (body as any).data?.token
      const userId = (body as any).user_id || (body as any).data?.user_id
      if (token) {
        localStorage.setItem('token', token)
        onLogin(token, undefined)
      } else if (userId) {
        localStorage.setItem('user_id', String(userId))
        onLogin(undefined, String(userId))
      } else {
        setError('Unexpected login response')
      }
    } catch (err: any) {
      setError(err?.body?.message || err?.body || err.message || 'Network error')
    }
  }

  if (showRegister) {
    return <Register onRegister={onLogin} onCancel={() => setShowRegister(false)} />
  }

  return (
    <div className="max-w-md mx-auto mt-20 p-6 bg-white rounded shadow">
      <h2 className="text-2xl font-bold mb-4">Login</h2>
      {error && <div className="text-sm text-red-600 mb-2">{error}</div>}
      <form onSubmit={submit} className="space-y-3">
        <input className="w-full p-2 border rounded" placeholder="phone or email" value={identifier} onChange={(e) => setIdentifier(e.target.value)} />
        <input className="w-full p-2 border rounded" placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <div className="flex justify-between items-center">
          <button className="px-4 py-2 bg-blue-600 text-white rounded" type="submit">Login</button>
          <button type="button" className="text-sm text-blue-600" onClick={() => setShowRegister(true)}>Register</button>
        </div>
      </form>
      <p className="text-xs text-gray-500 mt-4">Login will use the backend JWT handler; phone-first UX recommended.</p>
    </div>
  )
}
