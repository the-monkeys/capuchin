import { env } from "@/config/env"
import type { Todo } from "@/types"

const BASE = env.apiUrl

const authHeaders = (token: string): HeadersInit => ({
  "Content-Type": "application/json",
  Authorization: `Bearer ${token}`,
})

// Auth 

export const authApi = {
  login: async (email: string, password: string): Promise<{ token: string }> => {
    const res = await fetch(`${BASE}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })
    if (!res.ok) throw new Error("Invalid email or password")
    return res.json()
  },

  signup: async (email: string, password: string): Promise<void> => {
    const res = await fetch(`${BASE}/signup`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })
    if (!res.ok) {
      const data = await res.json()
      throw { response: { data } }
    }
  },
}

// Todo

const normalise = (raw: Record<string, unknown>): Todo => ({
  id: Number(raw.id ?? raw.ID),
  item: String(raw.item ?? raw.Item ?? ""),
  completed: Boolean(raw.completed ?? raw.Completed ?? false),
})

export const todosApi = {
  getAll: async (token: string): Promise<Todo[]> => {
    const res = await fetch(`${BASE}/api/user/todo`, {
      headers: authHeaders(token),
    })
    if (!res.ok) throw new Error(`Failed to fetch todos: ${res.status}`)
    const data = await res.json()
    return (data ?? []).map(normalise)
  },

  create: async (token: string, item: string): Promise<Todo> => {
    const res = await fetch(`${BASE}/api/user/todo`, {
      method: "POST",
      headers: authHeaders(token),
      body: JSON.stringify({ item, completed: false }),
    })
    if (!res.ok) throw new Error(`Failed to create todo: ${res.status}`)
    return normalise(await res.json())
  },

  toggle: async (token: string, id: number, completed: boolean): Promise<void> => {
    const res = await fetch(`${BASE}/api/user/todo/${id}`, {
      method: "PATCH",
      headers: authHeaders(token),
      body: JSON.stringify({ completed }),
    })
    if (!res.ok) throw new Error(`Failed to toggle todo: ${res.status}`)
  },

  update: async (token: string, id: number, item: string): Promise<Todo> => {
    const res = await fetch(`${BASE}/api/user/todo/${id}`, {
      method: "PATCH",
      headers: authHeaders(token),
      body: JSON.stringify({ item }),
    })
    if (!res.ok) throw new Error(`Failed to update todo: ${res.status}`)
    return normalise(await res.json())
  },

  delete: async (token: string, id: number): Promise<void> => {
    const res = await fetch(`${BASE}/api/user/todo/${id}`, {
      method: "DELETE",
      headers: authHeaders(token),
    })
    if (!res.ok) throw new Error(`Failed to delete todo: ${res.status}`)
  },
}