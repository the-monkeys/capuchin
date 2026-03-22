import { type ButtonHTMLAttributes } from "react"
import { Loader2 } from "lucide-react"

type Variant = "primary" | "secondary" | "ghost" | "danger"
type Size = "sm" | "md" | "lg"

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  size?: Size
  fullWidth?: boolean
  loading?: boolean
}

const variantClasses: Record<Variant, string> = {
  primary:
    "bg-[#111] text-white border-[1.5px] border-[#111] hover:bg-[#333] hover:-translate-y-px disabled:opacity-55 disabled:cursor-not-allowed",
  secondary:
    "bg-[#f9fafb] text-[#374151] border-[1.5px] border-[#e5e7eb] hover:bg-[#f3f4f6] hover:border-[#d1d5db] hover:text-[#111]",
  ghost:
    "bg-transparent text-[#374151] border-[1.5px] border-[#e5e7eb] hover:border-[#111] hover:text-[#111]",
  danger:
    "bg-transparent text-[#374151] border-[1.5px] border-[#e5e7eb] hover:bg-[#fef2f2] hover:text-[#ef4444] hover:border-[#fecaca]",
}

const sizeClasses: Record<Size, string> = {
  sm: "py-[5px] px-[13px] text-[0.77rem]",
  md: "py-2 px-[18px] text-[0.83rem]",
  lg: "py-[13px] px-6 text-[0.9rem] font-bold",
}

export function Button({ variant = "primary", size = "md", fullWidth = false, loading = false, className = "", children, ...props }: ButtonProps) {
  return (
    <button
      disabled={loading || props.disabled}
      className={[
        "rounded-xl font-semibold cursor-pointer transition-all duration-150 flex items-center justify-center gap-1.5",
        variantClasses[variant],
        sizeClasses[size],
        fullWidth ? "w-full" : "",
        className,
      ].join(" ")}
      {...props}
    >
      {loading ? <Loader2 size={16} className="animate-spin" /> : children}
    </button>
  )
}