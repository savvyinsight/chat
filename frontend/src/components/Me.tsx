import React, { useEffect, useState } from 'react'
import api from '../api'

export default function Me({ onProfileUpdate }: { onProfileUpdate?: (u: any) => void }) {
  const [profile, setProfile] = useState<any>(null)
  const [name, setName] = useState('')
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    api.get('/user/me').then((b) => {
      setProfile(b.data)
      setName(b.data.Name || b.data.name || '')
      setAvatarPreview(b.data.AvatarURL || b.data.avatar_url || null)
      if (onProfileUpdate) onProfileUpdate(b.data)
    }).catch(() => {})
  }, [])

  const upload = async (file: File | null) => {
    if (!file) return
    setLoading(true)
    const fd = new FormData()
    fd.append('avatar', file)
    const token = localStorage.getItem('token')
    const headers: any = {}
    if (token) headers['Authorization'] = `Bearer ${token}`
    try {
      const res = await fetch('/user/avatar', { method: 'POST', body: fd, headers })
      const body = await res.json()
      if (res.ok) {
        setAvatarPreview(body.avatar_url)
        // refresh profile
        const p = await api.get('/user/me')
        setProfile(p.data)
        if (onProfileUpdate) onProfileUpdate(p.data)
      }
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  const save = async () => {
    if (!profile) return
    const payload: any = { name }
    try {
      await api.put(`/user/${profile.ID}`, payload)
      const p = await api.get('/user/me')
      setProfile(p.data)
      if (onProfileUpdate) onProfileUpdate(p.data)
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <div className="max-w-md mx-auto p-4 bg-white rounded shadow">
      <h2 className="text-xl font-semibold mb-3">Me</h2>
      {profile ? (
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <img src={avatarPreview || '/static/avatars/default.png'} alt="avatar" className="w-16 h-16 rounded-full object-cover" />
            <div>
              <div className="text-sm text-gray-500">ID: {profile.ID}</div>
              <div className="text-lg font-medium">{profile.Name || profile.name}</div>
            </div>
          </div>

          <div>
            <label className="block text-sm">Display name</label>
            <input className="w-full p-2 border rounded" value={name} onChange={(e) => setName(e.target.value)} />
            <div className="flex gap-2 mt-2">
              <button className="px-3 py-1 bg-blue-600 text-white rounded" onClick={save}>Save</button>
            </div>
          </div>

          <div>
            <label className="block text-sm">Change avatar</label>
            <input type="file" accept="image/*" onChange={(e) => upload(e.target.files ? e.target.files[0] : null)} />
            {loading && <div className="text-sm text-gray-500">Uploading...</div>}
          </div>
        </div>
      ) : (
        <div>Loading...</div>
      )}
    </div>
  )
}
