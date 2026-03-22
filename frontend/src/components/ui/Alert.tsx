import { AlertTriangle, CheckCircle2 } from "lucide-react"

interface AlertProps {
  type: "error" | "success"
  message: string
}

const styles: Record<AlertProps["type"], string> = {
  error: "bg-[#fef2f2] border-[#fecaca] text-[#dc2626]",
  success: "bg-[#f0fdf4] border-[#bbf7d0] text-[#16a34a]",
}

export function Alert({ type, message }: AlertProps) {
  return (
    <div className={`mt-3 border-[1.5px] rounded-[10px] py-2.5 px-3.5 text-[0.82rem] font-medium flex items-center gap-2 ${styles[type]}`}>
      {type === "error" ? <AlertTriangle size={14} /> : <CheckCircle2 size={14} />}
      {message}
    </div>
  )
}