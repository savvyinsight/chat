export type Message = {
  type: string
  from?: number
  to?: number
  id?: number
  body?: string
}

export function createSocket(opts: { token?: string; userId?: number }) {
  const protocol = location.protocol === 'https:' ? 'wss' : 'ws'
  let url = `${protocol}://${location.hostname}:8080/ws`
  if (opts.token) url += `?token=${encodeURIComponent(opts.token)}`
  else if (opts.userId) url += `?user_id=${opts.userId}`
  const ws = new WebSocket(url)
  return ws
}
