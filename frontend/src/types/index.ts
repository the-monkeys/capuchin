export interface Todo {
  id: number | string
  item: string
  completed: boolean
}

export type FilterType = "all" | "active" | "done"

export interface AuthPayload {
  email: string
  exp: number
}