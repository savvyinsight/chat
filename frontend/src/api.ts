type Body = Record<string, any>

function getToken(): string | null {
  return localStorage.getItem('token')
}

async function request(path: string, opts: RequestInit = {}) {
  const headers = new Headers(opts.headers || undefined)
  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)
  const res = await fetch(path, { ...opts, headers })
  const text = await res.text()
  let body: any = {}
  try {
    body = text ? JSON.parse(text) : {}
  } catch (e) {
    body = { raw: text }
  }
  if (!res.ok) throw { status: res.status, body }
  return body
}

export async function apiGet(path: string) {
  return request(path, { method: 'GET' })
}

export async function apiPost(path: string, body: Body) {
  return request(path, { method: 'POST', body: JSON.stringify(body), headers: { 'Content-Type': 'application/json' } })
}

export default { get: apiGet, post: apiPost }
