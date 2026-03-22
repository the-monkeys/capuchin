import { useState } from "react"
import { useNavigate } from "react-router"
import { ArrowRight, ArrowLeft } from "lucide-react"
import { Navbar } from "@/components/layout/Navbar"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Alert } from "@/components/ui/Alert"
import { authApi } from "@/lib/api"

export default function Signup() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSignup = async () => {
    if (!email || !password) return
    setLoading(true)
    setError("")
    setSuccess("")
    try {
      await authApi.signup(email, password)
      setSuccess("Account created! Taking you to login.")
      setTimeout(() => navigate("/login"), 1500)
    } catch (err: unknown) {
      const anyErr = err as { response?: { data?: { error?: string } } }
      setError(anyErr.response?.data?.error ?? "Signup failed. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  const onKey = (e: React.KeyboardEvent) => e.key === "Enter" && handleSignup()

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
            Create your account
          </h1>

          <div className="flex flex-col gap-4 mt-8">
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
            onClick={handleSignup}
            loading={loading}
          >
            Create account <ArrowRight size={16} />
          </Button>

          {error && <Alert type="error" message={error} />}
          {success && <Alert type="success" message={success} />}

          <div className="flex items-center gap-3 my-6">
            <div className="flex-1 h-px bg-[#f3f4f6]" />
            <span className="text-[#d1d5db] font-semibold text-[0.75rem]">already have an account?</span>
            <div className="flex-1 h-px bg-[#f3f4f6]" />
          </div>

          <Button
            variant="secondary"
            fullWidth
            className="rounded-xl"
            onClick={() => navigate("/login")}
          >
            Sign in instead
          </Button>
        </div>
      </div>
    </div>
  )
}