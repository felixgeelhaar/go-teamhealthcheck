import { useState, useEffect, useCallback } from 'react'

export function useApi<T>(url: string | null) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(false)

  const refetch = useCallback(async () => {
    if (!url) return
    setLoading(true)
    try {
      const res = await fetch(url)
      if (res.ok) {
        setData(await res.json())
      }
    } finally {
      setLoading(false)
    }
  }, [url])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { data, loading, refetch }
}
