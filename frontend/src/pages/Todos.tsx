import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router'
import CapuchinSvg from '../assets/capuchin.svg'

interface Todo {
  id: string;
  item: string;
  completed: boolean;
}
//this key is used to store todos for unauthenticated users in localStorage. 
const LOCAL_STORAGE_KEY = "capuchin_guest_todos"
//uuid generator
const genId = () => crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).slice(2)
// capuchin logo
const MonkeyLogo = () => (
  <img src={CapuchinSvg} alt="Capuchin" style={{ width: '32px', height: '32px' }} />
)

const CheckIcon = () => (
  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="20 6 9 17 4 12" />
  </svg>
)
//load todos for unauthenticated users from localStorage, with error handling
const loadGuestTodos = (): Todo[] => {
  try {
    const raw = localStorage.getItem(LOCAL_STORAGE_KEY)
    return raw ? JSON.parse(raw) : []
  } catch { return [] }
}
//save todos for unauthenticated users to localStorage
const saveGuestTodos = (todos: Todo[]) => {
  localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(todos))
}
//main component
export default function Todos() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [inputValue, setInputValue] = useState("")
  const [loading, setLoading] = useState(true)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editValue, setEditValue] = useState("")
  const [filter, setFilter] = useState<"all" | "active" | "done">("all")
  const navigate = useNavigate()

  const token = localStorage.getItem("token")
  const isAuthed = Boolean(token)
  const API_BASE = "http://localhost:8080/user/todos"

  const userEmail = (() => {
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split(".")[1]))
      return payload.email ?? null
    } catch { return null }
  })()

  const authHeaders = (): HeadersInit => ({
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  })

  const normalise = (raw: any): Todo => ({
    id: String(raw.id ?? raw.ID),
    item: raw.item ?? raw.Item ?? "",
    completed: raw.completed ?? raw.Completed ?? false,
  })

  const handleLogout = () => {
    localStorage.removeItem("token")
    navigate("/login")
  }

  const updateGuest = (updater: (prev: Todo[]) => Todo[]) => {
    setTodos(prev => {
      const next = updater(prev)
      saveGuestTodos(next)
      return next
    })
  }

  useEffect(() => {
    if (!isAuthed) {
      setTodos(loadGuestTodos())
      setLoading(false)
      return
    }
    fetch(API_BASE, { headers: authHeaders() })
      .then(res => res.json())
      .then(data => { setTodos((data || []).map(normalise)); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])
//add a new todo
  const addTodo = async () => {
    if (!inputValue.trim()) return

    if (!isAuthed) {
      const newTodo: Todo = { id: genId(), item: inputValue.trim(), completed: false }
      updateGuest(prev => [...prev, newTodo])
      setInputValue("")
      return
    }

    try {
      const res = await fetch(API_BASE, {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify({ item: inputValue, completed: false })
      })
      const data = await res.json()
      setTodos(prev => [...prev, normalise(data)])
      setInputValue("")
    } catch (e) { console.error(e) }
  }
//toggle todo
  const toggleTodo = async (id: string, current: boolean) => {
    if (!isAuthed) {
      updateGuest(prev => prev.map(t => t.id === id ? { ...t, completed: !t.completed } : t))
      return
    }
    setTodos(prev => prev.map(t => t.id === id ? { ...t, completed: !t.completed } : t))
    try {
      await fetch(`${API_BASE}/${id}`, { method: "PATCH", headers: authHeaders() })
    } catch {
      setTodos(prev => prev.map(t => t.id === id ? { ...t, completed: current } : t))
    }
  }
//delete todo
  const deleteTodo = async (id: string) => {
    if (!isAuthed) {
      updateGuest(prev => prev.filter(t => t.id !== id))
      return
    }
    const old = [...todos]
    setTodos(prev => prev.filter(t => t.id !== id))
    try { await fetch(`${API_BASE}/${id}`, { method: "DELETE", headers: authHeaders() }) }
    catch { setTodos(old) }
  }
//save edited todo
  const saveEdit = async (id: string) => {
    const val = editValue.trim()
    if (!val) return

    if (!isAuthed) {
      updateGuest(prev => prev.map(t => t.id === id ? { ...t, item: val } : t))
      setEditValue("")
      return
    }

    try {
      const res = await fetch(`${API_BASE}/${id}/edit`, {
        method: "PATCH",
        headers: authHeaders(),
        body: JSON.stringify({ item: val })
      })
      const updated = await res.json()
      setTodos(prev => prev.map(t => t.id === id ? normalise(updated) : t))
      setEditValue("")
    } catch (e) { console.error(e) }
  }

  const done = todos.filter(t => t.completed).length
  const active = todos.length - done
  const filtered = todos.filter(t =>
    filter === "all" ? true : filter === "done" ? t.completed : !t.completed
  )

  return (
    <div className="grid-bg min-h-screen">

      {/* Nav */}
      <nav className="bg-white/90 backdrop-blur-md border-b border-[#e5e7eb] sticky top-0 z-[100] px-12 flex items-center justify-between h-16 max-[640px]:px-4">
        <div className="flex items-center gap-[9px] text-[1.1rem] font-extrabold text-[#111] tracking-[-0.03em] select-none cursor-pointer">
          <MonkeyLogo />
          Capuchin
        </div>
        <div className="flex items-center gap-2.5">
          {isAuthed ? (
            <>
              <span className="text-[0.82rem] font-semibold text-[#6b7280]">
                {userEmail ?? "Signed in"}
              </span>
              <button
                className="bg-transparent text-[#374151] border-[1.5px] border-[#e5e7eb] py-2 px-[18px] rounded-full text-[0.83rem] font-semibold cursor-pointer transition-all duration-150 hover:border-[#111] hover:text-[#111]"
                onClick={handleLogout}
              >
                Log out
              </button>
            </>
          ) : (
            <>
              <button
                className="bg-transparent text-[#374151] border-[1.5px] border-[#e5e7eb] py-2 px-[18px] rounded-full text-[0.83rem] font-semibold cursor-pointer transition-all duration-150 hover:border-[#111] hover:text-[#111]"
                onClick={() => navigate("/login")}
              >
                Log in
              </button>
              <button
                className="bg-[#111] text-white border-[1.5px] border-[#111] py-2 px-[18px] rounded-full text-[0.83rem] font-semibold cursor-pointer transition-[background,transform] duration-150 hover:bg-[#333] hover:-translate-y-px"
                onClick={() => navigate("/signup")}
              >
                Sign up free
              </button>
            </>
          )}
        </div>
      </nav>

      {/* Hero */}
      <div className="text-center pt-[72px] px-6 max-[640px]:pt-[52px] max-[640px]:px-4">
        <h1 className="text-[clamp(2.8rem,6vw,4.8rem)] font-black leading-[1.03] tracking-[-0.045em] text-[#111] mb-[18px]">
          Capuchin
        </h1>
        <p className="text-[1rem] text-[#6b7280] max-w-[400px] mx-auto mb-10 leading-[1.65]">
          handle, plan, and execute tasks
        </p>

        {/* Add bar */}
        <div className="flex items-center gap-2.5 max-w-[560px] mx-auto bg-white border-[1.5px] border-[#e5e7eb] rounded-full py-1.5 pl-[22px] pr-1.5 shadow-[0_4px_24px_rgba(0,0,0,0.07)] transition-[border-color,box-shadow] duration-200 focus-within:border-[#111] focus-within:shadow-[0_4px_32px_rgba(0,0,0,0.12)]">
          <input
            type="text"
            className="flex-1 border-none outline-none text-[0.92rem] text-[#111] bg-transparent placeholder:text-[#c4c9d4]"
            value={inputValue}
            onChange={e => setInputValue(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && addTodo()}
            placeholder="Add a task and press Enter..."
          />
          <button
            className="bg-[#111] text-white border-none rounded-full py-2.5 px-5 text-[0.85rem] font-bold cursor-pointer whitespace-nowrap transition-colors duration-150 hover:bg-[#333]"
            onClick={addTodo}
          >
            + Add task
          </button>
        </div>
      </div>

      {/* Task section */}
      <div className="max-w-[640px] mx-auto px-6 pb-20 mt-[52px] max-[640px]:px-4 max-[640px]:mt-12">

        {todos.length > 0 && (
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2.5">
              <span className="text-[0.72rem] font-bold text-[#9ca3af] uppercase tracking-[0.1em]">
                Your Tasks
              </span>
              <span className="text-[0.72rem] font-bold bg-[#f3f4f6] text-[#6b7280] rounded-full py-[3px] px-2.5">
                {active} left
              </span>
            </div>
            <div className="flex gap-1">
              {(["all", "active", "done"] as const).map(f => (
                <button
                  key={f}
                  className={`py-[5px] px-[13px] rounded-full text-[0.77rem] font-semibold cursor-pointer border-[1.5px] transition-all duration-150
                    ${filter === f
                      ? 'bg-[#111] text-white border-[#111]'
                      : 'bg-transparent text-[#9ca3af] border-transparent hover:text-[#374151]'
                    }`}
                  onClick={() => setFilter(f)}
                >
                  {f.charAt(0).toUpperCase() + f.slice(1)}
                </button>
              ))}
            </div>
          </div>
        )}

        {loading ? (
          <div className="text-center py-14 px-6 text-[#b0b7c3]">
            <div className="text-[2.6rem] mb-2.5 opacity-50">🐒</div>
            <div className="text-[0.88rem] font-medium">Loading your tasks...</div>
          </div>
        ) : todos.length > 0 && filtered.length === 0 ? (
          <div className="text-center py-14 px-6 text-[#b0b7c3]">
            <div className="text-[2.6rem] mb-2.5 opacity-50">🐒</div>
            <div className="text-[0.88rem] font-medium">No {filter} tasks</div>
          </div>
        ) : todos.length > 0 ? (
          <div className="flex flex-col gap-2">
            {filtered.map(todo => (
              <div
                key={todo.id}
                className={`border-[1.5px] rounded-[14px] py-[15px] px-[18px] flex items-center gap-[13px] shadow-[0_1px_4px_rgba(0,0,0,0.03)] transition-all duration-[180ms] ease-[ease] animate-slide-in hover:shadow-[0_4px_16px_rgba(0,0,0,0.06)] hover:-translate-y-px group
                  ${todo.completed
                    ? 'bg-[#fafafa] border-[#f5f5f5]'
                    : 'bg-white border-[#d1d5db] hover:border-[#adb5bd]'
                  }`}
              >
                {/* Checkbox */}
                <button
                  className={`w-[21px] h-[21px] rounded-full border-2 cursor-pointer flex items-center justify-center flex-shrink-0 transition-all duration-150
                    ${todo.completed
                      ? 'bg-[#111] border-[#111] text-white'
                      : 'bg-[#111] border-[#111] hover:bg-[#333] hover:border-[#333]'
                    }`}
                  onClick={() => toggleTodo(todo.id, todo.completed)}
                >
                  {todo.completed && <CheckIcon />}
                </button>

                {/* Content */}
                <div className="flex-1 min-w-0">
                  {editingId === todo.id ? (
                    <input
                      className="w-full border-[1.5px] border-[#111] rounded-lg py-[5px] px-2.5 text-[0.9rem] font-semibold text-[#111] outline-none bg-white"
                      value={editValue}
                      autoFocus
                      onChange={e => setEditValue(e.target.value)}
                      onKeyDown={e => {
                        if (e.key === "Enter") { setEditingId(null); saveEdit(todo.id) }
                        if (e.key === "Escape") setEditingId(null)
                      }}
                      onBlur={() => { setEditingId(null); saveEdit(todo.id) }}
                    />
                  ) : (
                    <span
                      className={`text-[0.9rem] cursor-pointer transition-colors duration-150 block
                        ${todo.completed
                          ? 'text-[#b0b7c3] line-through font-medium'
                          : 'font-semibold text-[#111]'
                        }`}
                      onClick={() => toggleTodo(todo.id, todo.completed)}
                    >
                      {todo.item}
                    </span>
                  )}
                </div>

                {/* Actions */}
                <div className="flex gap-0.5 opacity-0 transition-opacity duration-150 group-hover:opacity-100">
                  <button
                    className="w-7 h-7 border-none rounded-lg bg-transparent text-[#111] cursor-pointer flex items-center justify-center transition-all duration-150 hover:bg-[#eff6ff] hover:text-[#3b82f6]"
                    title="Edit"
                    onClick={() => { setEditingId(todo.id); setEditValue(todo.item) }}
                  >
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                    </svg>
                  </button>
                  <button
                    className="w-7 h-7 border-none rounded-lg bg-transparent text-[#111] cursor-pointer flex items-center justify-center transition-all duration-150 hover:bg-[#fef2f2] hover:text-[#ef4444]"
                    title="Delete"
                    onClick={() => deleteTodo(todo.id)}
                  >
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="3 6 5 6 21 6" />
                      <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                      <path d="M10 11v6M14 11v6" />
                      <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
                    </svg>
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : null}

      </div>
    </div>
  )
}