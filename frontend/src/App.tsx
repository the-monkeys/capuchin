import { Routes, Route } from "react-router"
import Login from "./pages/Login"
import Signup from "./pages/Signup"
import Todos from "./pages/Todos"

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Todos />} />
      <Route path="/todos" element={<Todos />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/login" element={<Login />} />
    </Routes>
  )
}