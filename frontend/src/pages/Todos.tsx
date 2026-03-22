import { useNavigate } from "react-router"
import { LogOut, LogIn, UserPlus } from "lucide-react"
import { Navbar } from "@/components/layout/Navbar"
import { Button } from "@/components/ui/Button"
import { EmptyState } from "@/components/ui/EmptyState"
import { TodoInput } from "@/components/todos/TodoInput"
import { TodoItem } from "@/components/todos/TodoItem"
import { TodoFilters } from "@/components/todos/TodoFilters"
import { useAuth } from "@/hooks/useAuth"
import { useTodos } from "@/hooks/useTodos"

export default function Todos() {
  const navigate = useNavigate()
  const { token, isAuthed, userEmail, logout } = useAuth()
  const { filtered, todos, loading, filter, setFilter, active, addTodo, toggleTodo, deleteTodo, updateTodo } = useTodos(token, isAuthed)

  const navRight = isAuthed ? (
    <>
      <span className="text-[0.82rem] font-semibold text-[#6b7280]">{userEmail ?? "Signed in"}</span>
      <Button variant="ghost" size="md" className="rounded-full" onClick={logout}>
        <LogOut size={14} />
        Log out
      </Button>
    </>
  ) : (
    <>
      <Button variant="ghost" size="md" className="rounded-full" onClick={() => navigate("/login")}>
        <LogIn size={14} />
        Log in
      </Button>
      <Button variant="primary" size="md" className="rounded-full" onClick={() => navigate("/signup")}>
        <UserPlus size={14} />
        Sign up free
      </Button>
    </>
  )

  return (
    <div className="grid-bg min-h-screen">
      <Navbar right={navRight} />

      <div className="text-center pt-18 px-6 max-[640px]:pt-13 max-[640px]:px-4">
        <h1 className="text-[clamp(2.8rem,6vw,4.8rem)] font-black leading-[1.03] tracking-[-0.045em] text-[#111] mb-4.5">
          Capuchin
        </h1>
        <p className="text-[1rem] text-[#6b7280] max-w-100 mx-auto mb-10 leading-[1.65]">
          handle, plan, and execute tasks
        </p>
        <TodoInput onAdd={addTodo} />
      </div>

      {/* Task list */}
      <div className="max-w-160 mx-auto px-6 pb-20 mt-13 max-[640px]:px-4 max-[640px]:mt-12">
        {todos.length > 0 && (
          <TodoFilters filter={filter} active={active} onChange={setFilter} />
        )}

        {loading ? (
          <EmptyState message="Loading your tasks..." />
        ) : todos.length === 0 ? null : filtered.length === 0 ? (
          <EmptyState message={`No ${filter} tasks`} />
        ) : (
          <div className="flex flex-col gap-2">
            {filtered.map((todo) => (
              <TodoItem
                key={todo.id}
                todo={todo}
                onToggle={toggleTodo}
                onDelete={deleteTodo}
                onUpdate={updateTodo}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}