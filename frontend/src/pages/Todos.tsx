import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router'

interface Todo {
  id: string;
  item: string;
  completed: boolean;
}

const MonkeyLogo = () => (
  <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
    <ellipse cx="5.5" cy="17" rx="4" ry="4.5" fill="#111"/>
    <ellipse cx="5.5" cy="17" rx="2.2" ry="2.8" fill="#c2855a"/>
    <ellipse cx="26.5" cy="17" rx="4" ry="4.5" fill="#111"/>
    <ellipse cx="26.5" cy="17" rx="2.2" ry="2.8" fill="#c2855a"/>
    <ellipse cx="16" cy="15" rx="11" ry="11" fill="#111"/>
    <ellipse cx="16" cy="18" rx="7" ry="6" fill="#c2855a"/>
    <circle cx="12.5" cy="13" r="2.2" fill="white"/>
    <circle cx="19.5" cy="13" r="2.2" fill="white"/>
    <circle cx="13" cy="13.4" r="1.1" fill="#111"/>
    <circle cx="20" cy="13.4" r="1.1" fill="#111"/>
    <circle cx="13.4" cy="13" r="0.4" fill="white"/>
    <circle cx="20.4" cy="13" r="0.4" fill="white"/>
    <ellipse cx="16" cy="16.5" rx="2" ry="1.4" fill="#8a5c3a"/>
    <circle cx="15.1" cy="16.3" r="0.5" fill="#5a3520"/>
    <circle cx="16.9" cy="16.3" r="0.5" fill="#5a3520"/>
    <path d="M13.5 19.5 Q16 21.5 18.5 19.5" stroke="#8a5c3a" strokeWidth="1.2" strokeLinecap="round" fill="none"/>
    <path d="M10 8 Q16 3 22 8" stroke="#111" strokeWidth="3" strokeLinecap="round" fill="none"/>
  </svg>
)

const CheckIcon = () => (
  <svg width="11" height="11" viewBox="0 0 12 12" fill="none">
    <path d="M2 6l3 3 5-5" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
)

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
  const API_URL = isAuthed
    ? "http://localhost:8080/user/todos"
    : "http://localhost:8080/todos"

  // Decode email from JWT payload 
  const userEmail = (() => {
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split(".")[1]))
      return payload.email ?? null
    } catch { return null }
  })()

  const authHeaders = (): HeadersInit => ({
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  })

  // Coerce server response — handles numeric id from Postgres and string id from JSON store
  const normalise = (raw: any): Todo => ({
    id: String(raw.id ?? raw.ID),
    item: raw.item ?? raw.Item ?? "",
    completed: raw.completed ?? raw.Completed ?? false,
  })

  const handleLogout = () => {
    localStorage.removeItem("token")
    navigate("/login")
  }

  useEffect(() => {
    fetch(API_URL, { headers: authHeaders() })
      .then(res => res.json())
      .then(data => { setTodos((data || []).map(normalise)); setLoading(false) })
      .catch(() => setLoading(false))
  }, [])
// add new todos
  const addTodo = async () => {
    if (!inputValue.trim()) return
    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify({ item: inputValue, completed: false })
      })
      const data = await res.json()
      setTodos(prev => [...prev, normalise(data)])
      setInputValue("")
    } catch (e) { console.error(e) }
  }
// toggle the todos
  const toggleTodo = async (id: string, current: boolean) => {
    setTodos(prev => prev.map(t => t.id === id ? { ...t, completed: !t.completed } : t))
    try {
      await fetch(`${API_URL}/${id}`, { method: "PATCH", headers: authHeaders() })
    } catch {
      setTodos(prev => prev.map(t => t.id === id ? { ...t, completed: current } : t))
    }
  }
// delete the todos
  const deleteTodo = async (id: string) => {
    const old = [...todos]
    setTodos(prev => prev.filter(t => t.id !== id))
    try { await fetch(`${API_URL}/${id}`, { method: "DELETE", headers: authHeaders() }) }
    catch { setTodos(old) }
  }
// edit the todos
  const saveEdit = async (id: string) => {
    const val = editValue.trim()
    if (!val) return
    try {
      const res = await fetch(`${API_URL}/${id}/edit`, {
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
    <> 
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Inter', sans-serif; background: #fafafa; color: #111; min-height: 100vh; }

        .page {
          min-height: 100vh;
          background-color: #fafafa;
          background-image:
            linear-gradient(rgba(0,0,0,0.055) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0,0,0,0.055) 1px, transparent 1px);
          background-size: 40px 40px;
        }

        .nav {
          background: rgba(255,255,255,0.9);
          backdrop-filter: blur(12px);
          border-bottom: 1px solid #e5e7eb;
          position: sticky; top: 0; z-index: 100;
          padding: 0 48px;
          display: flex; align-items: center; justify-content: space-between;
          height: 64px;
        }
        .nav-logo {
          display: flex; align-items: center; gap: 9px;
          font-size: 1.1rem; font-weight: 800; color: #111;
          letter-spacing: -0.03em; user-select: none; cursor: pointer;
        }
        .nav-right { display: flex; align-items: center; gap: 10px; }
        .btn-ghost {
          background: transparent; color: #374151;
          border: 1.5px solid #e5e7eb;
          padding: 8px 18px; border-radius: 999px;
          font-size: 0.83rem; font-weight: 600; cursor: pointer;
          font-family: 'Inter', sans-serif; transition: all 0.15s;
        }
        .btn-ghost:hover { border-color: #111; color: #111; }
        .btn-solid {
          background: #111; color: #fff; border: 1.5px solid #111;
          padding: 8px 18px; border-radius: 999px;
          font-size: 0.83rem; font-weight: 600; cursor: pointer;
          font-family: 'Inter', sans-serif; transition: background 0.15s, transform 0.1s;
        }
        .btn-solid:hover { background: #333; transform: translateY(-1px); }

        .hero {
          text-align: center;
          padding: 72px 24px 0;
        }
        .hero-title {
          font-size: clamp(2.8rem, 6vw, 4.8rem);
          font-weight: 900; line-height: 1.03;
          letter-spacing: -0.045em; color: #111; margin-bottom: 18px;
        }
        .hero-sub {
          font-size: 1rem; color: #6b7280;
          max-width: 400px; margin: 0 auto 40px; line-height: 1.65;
        }

        .add-bar {
          display: flex; align-items: center; gap: 10px;
          max-width: 560px; margin: 0 auto 0;
          background: #fff; border: 1.5px solid #e5e7eb;
          border-radius: 999px; padding: 6px 6px 6px 22px;
          box-shadow: 0 4px 24px rgba(0,0,0,0.07);
          transition: border-color 0.2s, box-shadow 0.2s;
        }
        .add-bar:focus-within {
          border-color: #111; box-shadow: 0 4px 28px rgba(0,0,0,0.12);
        }
        .add-bar input {
          flex: 1; border: none; outline: none;
          font-size: 0.9rem; font-family: 'Inter', sans-serif;
          color: #111; background: transparent;
        }
        .add-bar input::placeholder { color: #b0b7c3; }
        .add-bar-btn {
          background: #111; color: #fff; border: none;
          padding: 10px 22px; border-radius: 999px;
          font-size: 0.82rem; font-weight: 700;
          font-family: 'Inter', sans-serif; cursor: pointer;
          white-space: nowrap; transition: background 0.15s;
        }
        .add-bar-btn:hover { background: #333; }

        .task-section { max-width: 680px; margin: 64px auto 0; padding: 0 24px 80px; }
        .section-header {
          display: flex; align-items: center; justify-content: space-between;
          margin-bottom: 16px;
        }
        .section-title {
          font-size: 1.05rem; font-weight: 800;
          letter-spacing: -0.02em; color: #111;
        }
        .section-count {
          font-size: 0.78rem; color: #9ca3af; font-weight: 500;
          background: #f3f4f6; padding: 3px 10px; border-radius: 999px;
        }
        .filter-tabs {
          display: flex; gap: 4px; background: #f3f4f6;
          border-radius: 999px; padding: 4px;
        }
        .filter-tab {
          padding: 5px 14px; border-radius: 999px; border: none;
          background: transparent; font-size: 0.78rem; font-weight: 600;
          color: #6b7280; cursor: pointer; font-family: 'Inter', sans-serif;
          transition: all 0.15s;
        }
        .filter-tab.active { background: #fff; color: #111; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }

        .task-list { display: flex; flex-direction: column; gap: 8px; }
        .task-card {
          background: #fff; border: 1.5px solid #efefef;
          border-radius: 14px; padding: 15px 18px;
          display: flex; align-items: center; gap: 13px;
          box-shadow: 0 1px 4px rgba(0,0,0,0.03);
          transition: all 0.18s ease; animation: slideIn 0.22s ease;
        }
        @keyframes slideIn {
          from { opacity: 0; transform: translateY(6px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        .task-card:hover {
          border-color: #d1d5db; box-shadow: 0 4px 16px rgba(0,0,0,0.06);
          transform: translateY(-1px);
        }
        .task-card.done { background: #fafafa; border-color: #f5f5f5; }

        .check-btn {
          width: 21px; height: 21px; border-radius: 50%;
          border: 2px solid #d1d5db; background: transparent;
          cursor: pointer; display: flex; align-items: center; justify-content: center;
          flex-shrink: 0; transition: all 0.15s;
        }
        .check-btn:hover { border-color: #111; }
        .check-btn.checked { background: #111; border-color: #111; }

        .task-body { flex: 1; min-width: 0; }
        .task-text {
          font-size: 0.9rem; font-weight: 600; color: #111;
          cursor: pointer; transition: color 0.15s; display: block;
        }
        .task-text.done { color: #b0b7c3; text-decoration: line-through; font-weight: 500; }

        .task-actions { display: flex; gap: 2px; opacity: 0; transition: opacity 0.15s; }
        .task-card:hover .task-actions { opacity: 1; }
        .action-btn {
          width: 28px; height: 28px; border: none; border-radius: 8px;
          background: transparent; color: #c4c9d4; cursor: pointer;
          display: flex; align-items: center; justify-content: center;
          font-size: 0.8rem; transition: all 0.15s;
        }
        .action-btn.edit:hover  { background: #eff6ff; color: #3b82f6; }
        .action-btn.delete:hover { background: #fef2f2; color: #ef4444; }

        .edit-input {
          width: 100%; border: 1.5px solid #111; border-radius: 8px;
          padding: 5px 10px; font-size: 0.9rem; font-family: 'Inter', sans-serif;
          font-weight: 600; color: #111; outline: none; background: #fff;
        }

        .empty-state { text-align: center; padding: 56px 24px; color: #b0b7c3; }
        .empty-icon { font-size: 2.6rem; margin-bottom: 10px; opacity: 0.5; }
        .empty-text { font-size: 0.88rem; font-weight: 500; }

        .footer {
          text-align: center; padding: 28px;
          font-size: 0.73rem; color: #d1d5db;
          border-top: 1px solid #f3f4f6;
          margin-top: 40px;
        }

        @media (max-width: 640px) {
          .nav { padding: 0 16px; }
          .hero { padding: 52px 16px 0; }
          .task-section { padding: 0 16px 60px; margin-top: 48px; }
        }
      `}</style>

      <div className="page">

        <nav className="nav">
          <div className="nav-logo">
            <MonkeyLogo />
            Capuchin
          </div>
          <div className="nav-right">
            {isAuthed ? (
              <>
                <span style={{ fontSize: '0.82rem', fontWeight: 600, color: '#6b7280' }}>
                  {userEmail ?? "Signed in"}
                </span>
                <button className="btn-ghost" onClick={handleLogout}>Log out</button>
              </>
            ) : (
              <>
                <button className="btn-ghost" onClick={() => navigate("/login")}>Log in</button>
                <button className="btn-solid" onClick={() => navigate("/signup")}>Sign up free</button>
              </>
            )}
          </div>
        </nav>

        <div className="hero">
          <h1 className="hero-title">
            Capuchin
          </h1>
          <p className="hero-sub">
            handle, plan, and execute tasks
          </p>

          <div className="add-bar">
            <input
              type="text"
              value={inputValue}
              onChange={e => setInputValue(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && addTodo()}
              placeholder="Add a task and press Enter..."
            />
            <button className="add-bar-btn" onClick={addTodo}>+ Add task</button>
          </div>
        </div>

        <div className="task-section">
          {todos.length > 0 && (
            <div className="section-header">
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <div className="section-title">Your Tasks</div>
                <span className="section-count">{active} left</span>
              </div>
              <div className="filter-tabs">
                {(["all", "active", "done"] as const).map(f => (
                  <button
                    key={f}
                    className={`filter-tab ${filter === f ? 'active' : ''}`}
                    onClick={() => setFilter(f)}
                  >
                    {f.charAt(0).toUpperCase() + f.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          )}

          {loading ? (
            <div className="empty-state">
              <div className="empty-icon">🐒</div>
              <div className="empty-text">Loading your tasks...</div>
            </div>
          ) : todos.length > 0 && filtered.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">🐒</div>
              <div className="empty-text">No {filter} tasks</div>
            </div>
          ) : todos.length > 0 ? (
            <div className="task-list">
              {filtered.map(todo => (
                <div key={todo.id} className={`task-card ${todo.completed ? 'done' : ''}`}>
                  <button
                    className={`check-btn ${todo.completed ? 'checked' : ''}`}
                    onClick={() => toggleTodo(todo.id, todo.completed)}
                  >
                    {todo.completed && <CheckIcon />}
                  </button>

                  <div className="task-body">
                    {editingId === todo.id ? (
                      <input
                        className="edit-input"
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
                        className={`task-text ${todo.completed ? 'done' : ''}`}
                        onClick={() => toggleTodo(todo.id, todo.completed)}
                      >
                        {todo.item}
                      </span>
                    )}
                  </div>

                  <div className="task-actions">
                    <button
                      className="action-btn edit"
                      title="Edit"
                      onClick={() => { setEditingId(todo.id); setEditValue(todo.item) }}
                    >✏</button>
                    <button
                      className="action-btn delete"
                      title="Delete"
                      onClick={() => deleteTodo(todo.id)}
                    >✕</button>
                  </div>
                </div>
              ))}
            </div>
          ) : null}
        </div>

        <div className="footer">
          {todos.length > 0
            ? `${active} task${active !== 1 ? 's' : ''} remaining · Capuchin`
            : "Capuchin · Get things done 🐒"
          }
        </div>
      </div>
    </>
  )
}