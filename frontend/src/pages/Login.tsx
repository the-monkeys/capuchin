import { useState, useEffect } from "react"
import { useNavigate } from "react-router"
import { Navbar } from "@/components/layout/Navbar"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Alert } from "@/components/ui/Alert"
import { useAuth } from "@/hooks/useAuth"
import { authApi } from "@/lib/api"
import { ArrowRight, ArrowLeft } from "lucide-react"


export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { saveToken } = useAuth()

  useEffect(() => {
    if (localStorage.getItem("token")) navigate("/todos")
  }, [navigate])

  const handleLogin = async () => {
    if (!email || !password) return
    setLoading(true)
    setError("")
    try {
      const data = await authApi.login(email, password)
      saveToken(data.token)
      navigate("/todos")
    } catch {
      setError("Invalid email or password. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  const onKey = (e: React.KeyboardEvent) => e.key === "Enter" && handleLogin()

  return (
    <div className="grid-bg min-h-screen flex flex-col">
      <Navbar
        right={
          <Button
            variant="ghost"
            size="md"
            className="rounded-full"
            onClick={() => navigate("/")}
          >
            <ArrowLeft size={14} />
            Back to home
          </Button>
        }
      />

      <div className="flex-1 flex items-center justify-center p-10 px-6">
        <div className="bg-white border-[1.5px] border-[#e5e7eb] rounded-3xl p-12 px-11 w-full max-w-105 shadow-[0_8px_48px_rgba(0,0,0,0.06)] max-[480px]:p-9 max-[480px]:px-6">
          <h1 className="text-[1.85rem] font-black tracking-[-0.04em] text-[#111] mb-1.5 leading-[1.1]">
            Welcome back
          </h1>
          <p className="text-[0.875rem] text-[#9ca3af] mb-9 font-normal">
            Sign in to pick up where you left off
          </p>

          <div className="flex flex-col gap-4">
            <Input
              id="email"
              label="Email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={onKey}
            />
            <Input
              id="password"
              label="Password"
              type="password"
              placeholder="Enter password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyDown={onKey}
            />
          </div>

          <Button
            variant="primary"
            size="lg"
            fullWidth
            className="mt-6 rounded-xl"
            onClick={handleLogin}
            loading={loading}
          >
            Sign in <ArrowRight size={16} />
          </Button>

          {error && <Alert type="error" message={error} />}

          <div className="flex items-center gap-3 my-6">
            <div className="flex-1 h-px bg-[#f3f4f6]" />
            <span className="text-[#d1d5db] font-semibold text-[0.75rem]">or</span>
            <div className="flex-1 h-px bg-[#f3f4f6]" />
          </div>

          <Button
            variant="secondary"
            fullWidth
            className="rounded-xl"
            onClick={() => navigate("/signup")}
          >
            Don't have an account? Sign up
          </Button>
        </div>
      </div>
    </div>
  )
}