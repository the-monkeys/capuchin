import { type InputHTMLAttributes } from "react"

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export function Input({ label, id, className = "", ...props }: InputProps) {
  return (
    <div className="flex flex-col gap-2">
      {label && (
        <label htmlFor={id} className="text-[0.72rem] font-bold text-[#374151] uppercase tracking-[0.08em]">
          {label}
        </label>
      )}
      <input
        id={id}
        className={[
          "w-full py-3 px-4 border-[1.5px] border-[#e5e7eb] rounded-xl text-[0.9rem] text-[#111]",
          "outline-none bg-[#fafafa] transition-[border-color,box-shadow] duration-180",
          "focus:border-[#111] focus:shadow-[0_0_0_3px_rgba(0,0,0,0.06)] focus:bg-white",
          "placeholder:text-[#d1d5db]",
          className,
        ].join(" ")}
        {...props}
      />
    </div>
  )
}