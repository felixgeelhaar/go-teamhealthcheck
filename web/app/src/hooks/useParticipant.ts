import { useState } from 'react'

const STORAGE_KEY = 'participant_name'

export function useParticipant() {
  const [name, setName] = useState(localStorage.getItem(STORAGE_KEY) || '')

  const save = (n: string) => {
    localStorage.setItem(STORAGE_KEY, n)
    setName(n)
  }

  const clear = () => {
    localStorage.removeItem(STORAGE_KEY)
    setName('')
  }

  return { name, save, clear, isSet: name.length > 0 }
}
