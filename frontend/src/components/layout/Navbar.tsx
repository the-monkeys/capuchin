import { useNavigate } from "react-router"
import CapuchinSvg from "@/assets/capuchin.svg"

interface NavbarProps {
  right?: React.ReactNode
}

export function Navbar({ right }: NavbarProps) {
  const navigate = useNavigate()

  return (
    <nav className="bg-white/88 backdrop-blur-md border-b border-[#e5e7eb] h-16 px-12 flex items-center justify-between max-[480px]:px-5 sticky top-0 z-100">
      <button
        className="flex items-center gap-2.25 text-[1.1rem] font-extrabold text-[#111] tracking-[-0.03em] cursor-pointer select-none bg-transparent border-none"
        onClick={() => navigate("/")}
      >
        <img src={CapuchinSvg} alt="Capuchin" className="w-8 h-8" />
        Capuchin
      </button>
      {right && <div className="flex items-center gap-2.5">{right}</div>}
    </nav>
  )
}