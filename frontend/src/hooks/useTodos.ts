import { useState, useEffect, useCallback } from "react"
import { todosApi } from "@/lib/api"
import type { Todo, FilterType } from "@/types"

const GUEST_KEY = "capuchin_guest_todos"
const genId = () => crypto.randomUUID()

const loadGuestTodos = (): Todo[] => {
  try {
    const raw = localStorage.getItem(GUEST_KEY)
    return raw ? (JSON.parse(raw) as Todo[]) : []
  } catch {
    return []
  }
}

const saveGuestTodos = (todos: Todo[]) => {
  localStorage.setItem(GUEST_KEY, JSON.stringify(todos))
}

export function useTodos(token: string | null, isAuthed: boolean) {
  const [todos, setTodos] = useState<Todo[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<FilterType>("all")

  const updateGuest = useCallback((updater: (prev: Todo[]) => Todo[]) => {
    setTodos((prev) => {
      const next = updater(prev)
      saveGuestTodos(next)
      return next
    })
  }, [])

  useEffect(() => {
    if (!isAuthed) {
      setTodos(loadGuestTodos())
      setLoading(false)
      return
    }
    todosApi
      .getAll(token!)
      .then((data) => setTodos(data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [isAuthed, token])

  // CRUD operations

  const addTodo = useCallback(
    async (item: string) => {
      if (!item.trim()) return
      if (!isAuthed) {
        const newTodo: Todo = { id: genId(), item: item.trim(), completed: false }
        updateGuest((prev) => [...prev, newTodo])
        return
      }
      try {
        const created = await todosApi.create(token!, item)
        setTodos((prev) => [...prev, created])
      } catch (e) {
        console.error(e)
      }
    },
    [isAuthed, token, updateGuest],
  )

  const toggleTodo = useCallback(
    async (id: string, current: boolean) => {
      if (!isAuthed) {
        updateGuest((prev) => prev.map((t) => (t.id === id ? { ...t, completed: !t.completed } : t)))
        return
      }
      setTodos((prev) => prev.map((t) => (t.id === id ? { ...t, completed: !t.completed } : t)))
      try {
        await todosApi.toggle(token!, id, !current)
      } catch {
        setTodos((prev) => prev.map((t) => (t.id === id ? { ...t, completed: current } : t)))
      }
    },
    [isAuthed, token, updateGuest],
  )

  const deleteTodo = useCallback(
    async (id: string) => {
      if (!isAuthed) {
        updateGuest((prev) => prev.filter((t) => t.id !== id))
        return
      }
      const snapshot = todos
      setTodos((prev) => prev.filter((t) => t.id !== id))
      try {
        await todosApi.delete(token!, id)
      } catch {
        setTodos(snapshot)
      }
    },
    [isAuthed, token, todos, updateGuest],
  )

  const updateTodo = useCallback(
    async (id: string, item: string) => {
      if (!item.trim()) return
      if (!isAuthed) {
        updateGuest((prev) => prev.map((t) => (t.id === id ? { ...t, item: item.trim() } : t)))
        return
      }
      try {
        const updated = await todosApi.update(token!, id, item)
        setTodos((prev) => prev.map((t) => (t.id === id ? updated : t)))
      } catch (e) {
        console.error(e)
      }
    },
    [isAuthed, token, updateGuest],
  )
  
  const done = todos.filter((t) => t.completed).length
  const active = todos.length - done
  const filtered = todos.filter((t) => {
    if (filter === "all") return true
    if (filter === "done") return t.completed
    return !t.completed
  })

  return { todos, filtered, loading, filter, setFilter, active, done, addTodo, toggleTodo, deleteTodo, updateTodo }
}