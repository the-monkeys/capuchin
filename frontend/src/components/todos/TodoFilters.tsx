import type { FilterType } from "@/types"

interface TodoFiltersProps {
  filter: FilterType
  active: number
  onChange: (f: FilterType) => void
}

const FILTERS: FilterType[] = ["all", "active", "done"]

export function TodoFilters({ filter, active, onChange }: TodoFiltersProps) {
  return (
    <div className="flex items-center justify-between mb-4">
      <div className="flex items-center gap-2.5">
        <span className="text-[0.72rem] font-bold text-[#9ca3af] uppercase tracking-widest">Your Tasks</span>
        <span className="text-[0.72rem] font-bold bg-[#f3f4f6] text-[#6b7280] rounded-full py-0.75 px-2.5">
          {active} left
        </span>
      </div>
      <div className="flex gap-1">
        {FILTERS.map((f) => (
          <button
            key={f}
            onClick={() => onChange(f)}
            className={[
              "py-1.25 px-3.25 rounded-full text-[0.77rem] font-semibold cursor-pointer border-[1.5px] transition-all duration-150",
              filter === f
                ? "bg-[#111] text-white border-[#111]"
                : "bg-transparent text-[#9ca3af] border-transparent hover:text-[#374151]",
            ].join(" ")}
          >
            {f.charAt(0).toUpperCase() + f.slice(1)}
          </button>
        ))}
      </div>
    </div>
  )
}