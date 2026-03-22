import { useState, useCallback } from "react"
import { useNavigate } from "react-router"
import type { AuthPayload } from "@/types"

const TOKEN_KEY = "token"

export function useAuth() {
  const navigate = useNavigate()

  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem(TOKEN_KEY) 
  )

  const isAuthed = Boolean(token)

  const userEmail: string | null = (() => {
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split(".")[1])) as AuthPayload
      return payload.email ?? null
    } catch {
      return null
    }
  })()

  const saveToken = useCallback((t: string) => {
    localStorage.setItem(TOKEN_KEY, t)
    setToken(t)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
    navigate("/login")
  }, [navigate])

  const clearToken = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
  }, [])

  return { token, isAuthed, userEmail, saveToken, logout, clearToken }
}