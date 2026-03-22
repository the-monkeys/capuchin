import { useState } from "react"
import { Plus } from "lucide-react"

interface TodoInputProps {
  onAdd: (item: string) => void
}

export function TodoInput({ onAdd }: TodoInputProps) {
  const [value, setValue] = useState("")

  const handleAdd = () => {
    if (!value.trim()) return
    onAdd(value.trim())
    setValue("")
  }

  return (
    <div
      className="flex items-center gap-2.5 max-w-140 mx-auto bg-white border-[1.5px] border-[#e5e7eb] rounded-full py-1.5 pl-5.5 pr-1.5 shadow-[0_4px_24px_rgba(0,0,0,0.07)] transition-[border-color,box-shadow] duration-200 focus-within:border-[#111] focus-within:shadow-[0_4px_32px_rgba(0,0,0,0.12)]"
    >
      <input
        type="text"
        className="flex-1 border-none outline-none text-[0.92rem] text-[#111] bg-transparent placeholder:text-[#c4c9d4]"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleAdd()}
        placeholder="Add a task and press Enter..."
        aria-label="New task"
      />
      <button
        className="bg-[#111] text-white border-none rounded-full py-2.5 px-5 text-[0.85rem] font-bold cursor-pointer whitespace-nowrap transition-colors duration-150 hover:bg-[#333] flex items-center gap-1.5"
        onClick={handleAdd}
      >
        <Plus size={14} strokeWidth={2.5} />
        Add task
      </button>
    </div>
  )
}