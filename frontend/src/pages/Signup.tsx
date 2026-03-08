import { useState, useEffect } from "react"
import { signup } from "../api/auth"
import { useNavigate } from "react-router"
import CapuchinSvg from '../assets/capuchin.svg'

const MonkeyLogo = () => (
  <img src={CapuchinSvg} alt="Capuchin" style={{ width: '32px', height: '32px' }} />
)

export default function Signup() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    localStorage.removeItem("token")
  }, [])

  const handleSignup = async () => {
    if (!email || !password) return
    setLoading(true)
    try {
      setError("")
      setSuccess("")
      await signup(email, password)
      setSuccess("Account created! Taking you to login...")
      setTimeout(() => navigate("/login"), 1500)
    } catch (err: any) {
      setError(err.response?.data?.error || "Signup failed. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="grid-bg min-h-screen flex flex-col">
      {/* Nav */}
      <nav className="bg-white/[0.88] backdrop-blur-md border-b border-[#e5e7eb] h-16 px-12 flex items-center justify-between max-[480px]:px-5">
        <div
          className="flex items-center gap-[9px] text-[1.1rem] font-extrabold text-[#111] tracking-[-0.03em] cursor-pointer select-none"
          onClick={() => navigate("/")}
        >
          <MonkeyLogo />
          Capuchin
        </div>
        <button
          className="text-[0.83rem] font-semibold text-[#6b7280] cursor-pointer bg-transparent border-none flex items-center gap-1.5 transition-colors duration-150 hover:text-[#111]"
          onClick={() => navigate("/")}
        >
          ← Back to home
        </button>
      </nav>

      {/* Body */}
      <div className="flex-1 flex items-center justify-center p-10 px-6">
        <div className="bg-white border-[1.5px] border-[#e5e7eb] rounded-3xl p-12 px-11 w-full max-w-[420px] shadow-[0_8px_48px_rgba(0,0,0,0.06)] max-[480px]:p-9 max-[480px]:px-6">

          <span className="text-[2.8rem] leading-none mb-5 block">🐒</span>
          <h1 className="text-[1.85rem] font-black tracking-[-0.04em] text-[#111] mb-1.5 leading-[1.1]">
            Create your account
          </h1>

          {/* Email */}
          <div className="mb-4 mt-8">
            <label className="block text-[0.72rem] font-bold text-[#374151] uppercase tracking-[0.08em] mb-2">
              Email
            </label>
            <input
              className="w-full py-3 px-4 border-[1.5px] border-[#e5e7eb] rounded-xl text-[0.9rem] text-[#111] outline-none bg-[#fafafa] transition-[border-color,box-shadow] duration-[180ms] focus:border-[#111] focus:shadow-[0_0_0_3px_rgba(0,0,0,0.06)] focus:bg-white placeholder:text-[#d1d5db]"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={e => setEmail(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && handleSignup()}
            />
          </div>

          {/* Password */}
          <div className="mb-4">
            <label className="block text-[0.72rem] font-bold text-[#374151] uppercase tracking-[0.08em] mb-2">
              Password
            </label>
            <input
              className="w-full py-3 px-4 border-[1.5px] border-[#e5e7eb] rounded-xl text-[0.9rem] text-[#111] outline-none bg-[#fafafa] transition-[border-color,box-shadow] duration-[180ms] focus:border-[#111] focus:shadow-[0_0_0_3px_rgba(0,0,0,0.06)] focus:bg-white placeholder:text-[#d1d5db]"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={e => setPassword(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && handleSignup()}
            />
          </div>

          <button
            className="w-full mt-2 py-[13px] bg-[#111] text-white border-none rounded-xl text-[0.9rem] font-bold cursor-pointer transition-[background,transform] duration-150 hover:bg-[#333] hover:-translate-y-px disabled:opacity-[0.55] disabled:cursor-not-allowed"
            onClick={handleSignup}
            disabled={loading}
          >
            {loading ? "Creating account..." : "Create account →"}
          </button>

          {error && (
            <div className="mt-3 bg-[#fef2f2] border-[1.5px] border-[#fecaca] rounded-[10px] py-2.5 px-3.5 text-[0.82rem] text-[#dc2626] font-medium">
              ⚠ {error}
            </div>
          )}
          {success && (
            <div className="mt-3 bg-[#f0fdf4] border-[1.5px] border-[#bbf7d0] rounded-[10px] py-2.5 px-3.5 text-[0.82rem] text-[#16a34a] font-medium">
              ✓ {success}
            </div>
          )}

          {/* Divider */}
          <div className="flex items-center gap-3 my-6">
            <div className="flex-1 h-px bg-[#f3f4f6]" />
            <span className="text-[#d1d5db] font-semibold text-[0.75rem]">already have an account?</span>
            <div className="flex-1 h-px bg-[#f3f4f6]" />
          </div>

          <button
            className="w-full py-3 bg-[#f9fafb] text-[#374151] border-[1.5px] border-[#e5e7eb] rounded-xl text-[0.875rem] font-semibold cursor-pointer transition-all duration-150 hover:bg-[#f3f4f6] hover:border-[#d1d5db] hover:text-[#111]"
            onClick={() => navigate("/login")}
          >
            Sign in instead
          </button>

        </div>
      </div>
    </div>
  )
}