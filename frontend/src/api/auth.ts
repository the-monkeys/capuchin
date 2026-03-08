const BASE_URL = "http://localhost:8080"

export async function login(email: string, password: string) {
  const res = await fetch(`${BASE_URL}/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  })

  if (!res.ok) throw new Error("Login failed")

  return res.json()
}

export const signup = async (email: string, password: string) => {
  const res = await fetch("http://localhost:8080/signup", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ email, password })
  })

  if (!res.ok) {
    const data = await res.json()
    throw { response: { data } }
  }

  return res.json()
}


