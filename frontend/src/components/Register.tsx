import React, { useState } from 'react'
import api from '../api'

export default function Register({ onRegister, onCancel }: { onRegister: (token?: string, userId?: string) => void; onCancel?: () => void }) {
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [repassword, setRepassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    if (!password || password !== repassword) {
      setError('Passwords do not match')
      return
    }
    setLoading(true)
    try {
      const body = await api.post('/user/register', { name, phone, email, password, repassword })
      const token = (body as any).token || (body as any).data?.token
      const userId = (body as any).user_id || (body as any).data?.user_id
      if (token) {
        localStorage.setItem('token', token)
        onRegister(token, undefined)
      } else if (userId) {
        localStorage.setItem('user_id', String(userId))
        onRegister(undefined, String(userId))
      } else {
        setError('Unexpected register response')
      }
    } catch (err: any) {
      setError(err?.body?.message || err?.body || err.message || 'Network error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-12 p-6 bg-white rounded shadow">
      <h2 className="text-2xl font-bold mb-4">Register</h2>
      {error && <div className="text-sm text-red-600 mb-2">{error}</div>}
      <form onSubmit={submit} className="space-y-3">
        <input className="w-full p-2 border rounded" placeholder="display name" value={name} onChange={(e) => setName(e.target.value)} />
        <input className="w-full p-2 border rounded" placeholder="phone (preferred)" value={phone} onChange={(e) => setPhone(e.target.value)} />
        <input className="w-full p-2 border rounded" placeholder="email (optional)" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input className="w-full p-2 border rounded" placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <input className="w-full p-2 border rounded" placeholder="confirm password" type="password" value={repassword} onChange={(e) => setRepassword(e.target.value)} />
        <div className="flex justify-between items-center">
          <button className="px-4 py-2 bg-green-600 text-white rounded" type="submit" disabled={loading}>{loading ? 'Registering...' : 'Register'}</button>
          {onCancel && <button type="button" className="text-sm text-gray-600" onClick={onCancel}>Cancel</button>}
        </div>
      </form>
      <p className="text-xs text-gray-500 mt-4">We recommend registering with your phone number for easier login and recovery.</p>
    </div>
  )
}
