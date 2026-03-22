import { useState } from "react"
import { Check, Pencil, Trash2 } from "lucide-react"
import type { Todo } from "@/types"

interface TodoItemProps {
  todo: Todo
  onToggle: (id: string, current: boolean) => void
  onDelete: (id: string) => void
  onUpdate: (id: string, item: string) => void
}

export function TodoItem({ todo, onToggle, onDelete, onUpdate }: TodoItemProps) {
  const [editing, setEditing] = useState(false)
  const [editValue, setEditValue] = useState(todo.item)

  const commitEdit = () => {
    setEditing(false)
    if (editValue.trim()) onUpdate(todo.id, editValue.trim())
  }

  return (
    <div
      className={[
        "border-[1.5px] rounded-[14px] py-3.75 px-4.5 flex items-center gap-3.25",
        "shadow-[0_1px_4px_rgba(0,0,0,0.03)] transition-all duration-180 ease-[ease]",
        "animate-slide-in hover:shadow-[0_4px_16px_rgba(0,0,0,0.06)] hover:-translate-y-px group",
        todo.completed ? "bg-[#fafafa] border-[#f5f5f5]" : "bg-white border-[#d1d5db] hover:border-[#adb5bd]",
      ].join(" ")}
    >
      <button
        className={[
          "w-5.25 h-5.25 rounded-full border-2 cursor-pointer flex items-center justify-center shrink-0 transition-all duration-150",
          todo.completed
            ? "bg-[#111] border-[#111] text-white"
            : "bg-[#111] border-[#111] hover:bg-[#333] hover:border-[#333]",
        ].join(" ")}
        onClick={() => onToggle(todo.id, todo.completed)}
        aria-label={todo.completed ? "Mark as incomplete" : "Mark as complete"}
      >
        {todo.completed && <Check size={11} strokeWidth={3.5} />}
      </button>

      <div className="flex-1 min-w-0">
        {editing ? (
          <input
            autoFocus
            className="w-full border-[1.5px] border-[#111] rounded-lg py-1.25 px-2.5 text-[0.9rem] font-semibold text-[#111] outline-none bg-white"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") commitEdit()
              if (e.key === "Escape") setEditing(false)
            }}
            onBlur={commitEdit}
          />
        ) : (
          <span
            className={[
              "text-[0.9rem] cursor-pointer transition-colors duration-150 block truncate",
              todo.completed ? "text-[#b0b7c3] line-through font-medium" : "font-semibold text-[#111]",
            ].join(" ")}
            onClick={() => onToggle(todo.id, todo.completed)}
          >
            {todo.item}
          </span>
        )}
      </div>

      <div className="flex gap-0.5 opacity-0 transition-opacity duration-150 group-hover:opacity-100">
        <button
          aria-label="Edit task"
          className="w-7 h-7 border-none rounded-lg bg-transparent cursor-pointer flex items-center justify-center transition-all duration-150 text-[#6b7280] hover:bg-[#eff6ff] hover:text-[#3b82f6]"
          onClick={() => { setEditValue(todo.item); setEditing(true) }}
        >
          <Pencil size={13} strokeWidth={2.5} />
        </button>
        <button
          aria-label="Delete task"
          className="w-7 h-7 border-none rounded-lg bg-transparent cursor-pointer flex items-center justify-center transition-all duration-150 text-[#6b7280] hover:bg-[#fef2f2] hover:text-[#ef4444]"
          onClick={() => onDelete(todo.id)}
        >
          <Trash2 size={13} strokeWidth={2.5} />
        </button>
      </div>
    </div>
  )
}