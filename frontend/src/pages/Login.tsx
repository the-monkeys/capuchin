import { useState, useEffect } from "react"
import { login } from "../api/auth"
import { useNavigate } from "react-router"

const MonkeyLogo = () => (
  <svg width="34" height="34" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
    <ellipse cx="5.5" cy="17" rx="4" ry="4.5" fill="#111"/>
    <ellipse cx="5.5" cy="17" rx="2.2" ry="2.8" fill="#c2855a"/>
    <ellipse cx="26.5" cy="17" rx="4" ry="4.5" fill="#111"/>
    <ellipse cx="26.5" cy="17" rx="2.2" ry="2.8" fill="#c2855a"/>
    <ellipse cx="16" cy="15" rx="11" ry="11" fill="#111"/>
    <ellipse cx="16" cy="18" rx="7" ry="6" fill="#c2855a"/>
    <circle cx="12.5" cy="13" r="2.2" fill="white"/>
    <circle cx="19.5" cy="13" r="2.2" fill="white"/>
    <circle cx="13" cy="13.4" r="1.1" fill="#111"/>
    <circle cx="20" cy="13.4" r="1.1" fill="#111"/>
    <circle cx="13.4" cy="13" r="0.4" fill="white"/>
    <circle cx="20.4" cy="13" r="0.4" fill="white"/>
    <ellipse cx="16" cy="16.5" rx="2" ry="1.4" fill="#8a5c3a"/>
    <circle cx="15.1" cy="16.3" r="0.5" fill="#5a3520"/>
    <circle cx="16.9" cy="16.3" r="0.5" fill="#5a3520"/>
    <path d="M13.5 19.5 Q16 21.5 18.5 19.5" stroke="#8a5c3a" strokeWidth="1.2" strokeLinecap="round" fill="none"/>
    <path d="M10 8 Q16 3 22 8" stroke="#111" strokeWidth="3" strokeLinecap="round" fill="none"/>
  </svg>
)

export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem("token")
    if (token) navigate("/todos")
  }, [])

  const handleLogin = async () => {
    if (!email || !password) return
    setLoading(true)
    try {
      setError("")
      const data = await login(email, password)
      localStorage.setItem("token", data.token)
      navigate("/todos")
    } catch {
      setError("Invalid email or password. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Inter', sans-serif; background: #fafafa; color: #111; min-height: 100vh; }

        .page {
          min-height: 100vh;
          background-color: #fafafa;
          background-image:
            linear-gradient(rgba(0,0,0,0.055) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0,0,0,0.055) 1px, transparent 1px);
          background-size: 40px 40px;
          display: flex; flex-direction: column;
        }

        .nav {
          background: rgba(255,255,255,0.88);
          backdrop-filter: blur(12px);
          border-bottom: 1px solid #e5e7eb;
          height: 64px; padding: 0 48px;
          display: flex; align-items: center; justify-content: space-between;
        }
        .nav-logo {
          display: flex; align-items: center; gap: 9px;
          font-size: 1.1rem; font-weight: 800; color: #111;
          letter-spacing: -0.03em; cursor: pointer; user-select: none;
        }
        .nav-back {
          font-size: 0.83rem; font-weight: 600; color: #6b7280;
          cursor: pointer; background: none; border: none;
          font-family: 'Inter', sans-serif;
          display: flex; align-items: center; gap: 6px;
          transition: color 0.15s;
        }
        .nav-back:hover { color: #111; }

        .body {
          flex: 1; display: flex; align-items: center;
          justify-content: center; padding: 40px 24px;
        }

        .card {
          background: #fff;
          border: 1.5px solid #e5e7eb;
          border-radius: 24px;
          padding: 48px 44px;
          width: 100%; max-width: 420px;
          box-shadow: 0 8px 48px rgba(0,0,0,0.06);
        }

        .card-monkey {
          font-size: 2.8rem; line-height: 1;
          margin-bottom: 20px; display: block;
        }
        .card-title {
          font-size: 1.85rem; font-weight: 900;
          letter-spacing: -0.04em; color: #111;
          margin-bottom: 6px; line-height: 1.1;
        }
        .card-sub {
          font-size: 0.875rem; color: #9ca3af;
          margin-bottom: 36px; font-weight: 400;
        }

        .form-group { margin-bottom: 16px; }
        .form-label {
          display: block;
          font-size: 0.72rem; font-weight: 700; color: #374151;
          text-transform: uppercase; letter-spacing: 0.08em;
          margin-bottom: 8px;
        }
        .form-input {
          width: 100%; padding: 12px 16px;
          border: 1.5px solid #e5e7eb; border-radius: 12px;
          font-size: 0.9rem; font-family: 'Inter', sans-serif; color: #111;
          outline: none; background: #fafafa;
          transition: border-color 0.18s, box-shadow 0.18s;
        }
        .form-input:focus {
          border-color: #111;
          box-shadow: 0 0 0 3px rgba(0,0,0,0.06);
          background: #fff;
        }
        .form-input::placeholder { color: #d1d5db; }

        .submit-btn {
          width: 100%; margin-top: 8px;
          padding: 13px; background: #111; color: #fff;
          border: none; border-radius: 12px;
          font-size: 0.9rem; font-weight: 700;
          font-family: 'Inter', sans-serif; cursor: pointer;
          transition: background 0.15s, transform 0.1s;
        }
        .submit-btn:hover:not(:disabled) { background: #333; transform: translateY(-1px); }
        .submit-btn:disabled { opacity: 0.55; cursor: not-allowed; }

        .error-msg {
          margin-top: 12px;
          background: #fef2f2; border: 1.5px solid #fecaca;
          border-radius: 10px; padding: 10px 14px;
          font-size: 0.82rem; color: #dc2626; font-weight: 500;
        }

        .divider {
          display: flex; align-items: center; gap: 12px;
          margin: 24px 0; color: #e5e7eb; font-size: 0.75rem;
        }
        .divider::before, .divider::after {
          content: ''; flex: 1; height: 1px; background: #f3f4f6;
        }
        .divider span { color: #d1d5db; font-weight: 600; }

        .alt-btn {
          width: 100%; padding: 12px;
          background: #f9fafb; color: #374151;
          border: 1.5px solid #e5e7eb; border-radius: 12px;
          font-size: 0.875rem; font-weight: 600;
          font-family: 'Inter', sans-serif; cursor: pointer;
          transition: all 0.15s;
        }
        .alt-btn:hover { background: #f3f4f6; border-color: #d1d5db; color: #111; }

        @media (max-width: 480px) {
          .card { padding: 36px 24px; }
          .nav { padding: 0 20px; }
        }
      `}</style>

      <div className="page">
        <nav className="nav">
          <div className="nav-logo" onClick={() => navigate("/")}>
            <MonkeyLogo />
            Capuchin
          </div>
          <button className="nav-back" onClick={() => navigate("/")}>
            ← Back to home
          </button>
        </nav>

        <div className="body">
          <div className="card">
            <span className="card-monkey">🐒</span>
            <h1 className="card-title">Welcome back</h1>
            <p className="card-sub">Sign in to pick up where you left off</p>

            <div className="form-group">
              <label className="form-label">Email</label>
              <input
                className="form-input"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={e => setEmail(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && handleLogin()}
              />
            </div>

            <div className="form-group">
              <label className="form-label">Password</label>
              <input
                className="form-input"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={e => setPassword(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && handleLogin()}
              />
            </div>

            <button className="submit-btn" onClick={handleLogin} disabled={loading}>
              {loading ? "Signing in..." : "Sign in →"}
            </button>

            {error && <div className="error-msg">⚠ {error}</div>}

            <div className="divider"><span>or</span></div>

            <button className="alt-btn" onClick={() => navigate("/signup")}>
              Don't have an account? Sign up
            </button>
          </div>
        </div>
      </div>
    </>
  )
}