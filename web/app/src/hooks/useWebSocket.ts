import { useEffect, useRef, useCallback } from 'react'
import type { WSEvent } from '../types'

export function useWebSocket(onEvent: (event: WSEvent) => void) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimer = useRef<ReturnType<typeof setTimeout>>(undefined)
  const onEventRef = useRef<(event: WSEvent) => void>(onEvent)
  onEventRef.current = onEvent

  const connect = useCallback(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`)

    ws.onmessage = (e) => {
      try {
        const event: WSEvent = JSON.parse(e.data)
        onEventRef.current(event)
      } catch {
        // ignore malformed messages
      }
    }

    ws.onclose = () => {
      wsRef.current = null
      reconnectTimer.current = setTimeout(connect, 2000)
    }

    ws.onerror = () => ws.close()

    wsRef.current = ws
  }, [])

  useEffect(() => {
    connect()
    return () => {
      clearTimeout(reconnectTimer.current)
      wsRef.current?.close()
    }
  }, [connect])
}
