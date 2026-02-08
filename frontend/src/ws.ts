export type Message = {
  type: string
  from?: number
  to?: number
  id?: number
  body?: string
}

export function createSocket(opts: { token?: string; userId?: number }) {
  const protocol = location.protocol === 'https:' ? 'wss' : 'ws'
  // use explicit backend port 8080 (same as server)
  let url = `${protocol}://${location.hostname}:8080/ws`
  if (opts.token) url += `?token=${encodeURIComponent(opts.token)}`
  else if (opts.userId) url += `?user_id=${opts.userId}`

  // WebSocket wrapper with reconnection, backoff and outgoing buffer.
  let ws: WebSocket | null = null
  let reconnectAttempts = 0
  let reconnectTimer: number | null = null
  const sendBuffer: Array<string> = []

  // public-facing handlers (assignable)
  const wrapper: any = {
    onopen: null,
    onmessage: null,
    onclose: null,
    onerror: null,
    connected: false,
    send(data: string) {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data)
      } else {
        // buffer messages while disconnected
        sendBuffer.push(data)
      }
    },
    close() {
      // stop reconnect attempts and close underlying socket
      if (reconnectTimer) {
        clearTimeout(reconnectTimer)
        reconnectTimer = null
      }
      if (ws) {
        ws.close()
        ws = null
      }
    },
  }

  function connect() {
    ws = new WebSocket(url)
    ws.onopen = (ev) => {
      reconnectAttempts = 0
      wrapper.connected = true
      // flush buffer
      while (sendBuffer.length > 0) {
        const d = sendBuffer.shift()!
        try {
          ws!.send(d)
        } catch (e) {
          // if send fails, re-buffer and break
          sendBuffer.unshift(d)
          break
        }
      }
      if (typeof wrapper.onopen === 'function') wrapper.onopen(ev)
    }
    ws.onmessage = (ev) => {
      if (typeof wrapper.onmessage === 'function') wrapper.onmessage(ev)
    }
    ws.onclose = (ev) => {
      wrapper.connected = false
      if (typeof wrapper.onclose === 'function') wrapper.onclose(ev)
      // attempt reconnect unless closed explicitly
      scheduleReconnect()
    }
    ws.onerror = (ev) => {
      if (typeof wrapper.onerror === 'function') wrapper.onerror(ev)
    }
  }

  function scheduleReconnect() {
    if (reconnectTimer) return
    reconnectAttempts++
    // exponential backoff: 1s,2s,4s,8s,... cap 30s
    const delay = Math.min(30000, Math.pow(2, Math.min(reconnectAttempts, 8)) * 1000)
    reconnectTimer = window.setTimeout(() => {
      reconnectTimer = null
      connect()
    }, delay)
  }

  // start
  connect()

  return wrapper
}
