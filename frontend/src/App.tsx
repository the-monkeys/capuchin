import { useState, useEffect } from 'react'
import './App.css'

interface Todo {
  id: string;
  item: string;
  completed: boolean;
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [inputValue, setInputValue] = useState("")
  const [loading, setLoading] = useState(true)

  const API_URL = "http://localhost:8080/todos"

  useEffect(() => {
    fetch(API_URL)
      .then(res => res.json())
      .then(data => {
        setTodos(data || [])
        setLoading(false)
      })
      .catch(err => {
        console.error("Error fetching:", err)
        setLoading(false)
      })
  }, [])


  // Add a new todo
  const addTodo = async () => {
    if (inputValue.trim() === "") return

    const newTodo = { item: inputValue, completed: false }

    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newTodo)
      })

      const data = await res.json()
      setTodos([...todos, data])
      setInputValue("")
    } catch (error) {
      console.error(error)
    }
  }


  // Toggle a todo
  const toggleTodo = async (id: string, currentStatus: boolean) => {
    console.log("Toggling todo:", id, currentStatus)

    setTodos(todos.map(t =>
      t.id === id ? { ...t, completed: !t.completed } : t
    ))

    try {
      await fetch(`${API_URL}/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
      })
      console.log("Successfully toggled todo:", id)
    } catch (error) {
      console.error(error)


      setTodos(todos.map(t =>
        t.id === id ? { ...t, completed: currentStatus } : t
      ))
      console.log("Reverted todo:", id)
    }
  }


  // Delete a todo
  const deleteTodo = async (id: string) => {
    const oldTodos = [...todos]
    setTodos(todos.filter(t => t.id !== id))

    try {
      await fetch(`${API_URL}/${id}`, { method: "DELETE" })
    } catch (error) {
      console.error(error)
      setTodos(oldTodos)
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-gray-800 rounded-xl shadow-2xl overflow-hidden border border-gray-700">

        {/* Header */}
        <div className="p-6 bg-gray-800 border-b border-gray-700">
          <h1 className="text-3xl font-bold text-white tracking-tight">Capuchin Todos</h1>
          <p className="text-gray-400 text-sm mt-1">Manage your tasks efficiently</p>
        </div>

        {/* Input Area */}
        <div className="p-6 bg-gray-800/50">
          <div className="relative">
            <input
              type="text"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && addTodo()}
              placeholder="What needs doing?"
              className="w-full pl-4 pr-20 py-3 bg-gray-900 text-white rounded-lg border border-gray-700 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all placeholder-gray-500"
            />
            <button
              onClick={addTodo}
              className="absolute right-1.5 top-1.5 bottom-1.5 px-4 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-md transition-colors duration-200"
            >
              Add
            </button>
          </div>
        </div>

        {/* List Area */}
        <div className="bg-gray-800">
          {loading ? (
            <div className="p-8 text-center text-gray-400">Loading tasks...</div>
          ) : todos.length === 0 ? (
            <div className="p-8 text-center text-gray-500 italic">No tasks yet. Add one above!</div>
          ) : (
            <ul className="divide-y divide-gray-700">
              {todos.map((todo) => (
                <li key={todo.id} className="group flex items-center justify-between p-4 hover:bg-gray-700/30 transition-colors duration-150">
                  <div className="flex items-center gap-3 flex-1 overflow-hidden">
                    <button
                      onClick={() => toggleTodo(todo.id, todo.completed)}
                      className={`w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all duration-200 ${todo.completed
                        ? 'bg-green-500 border-green-500'
                        : 'border-gray-500 hover:border-gray-400'
                        }`}
                    >
                      {todo.completed && (
                        <svg className="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3.5" d="M5 13l4 4L19 7" />
                        </svg>
                      )}
                    </button>
                    <span
                      onClick={() => toggleTodo(todo.id, todo.completed)}
                      className={`text-lg cursor-pointer select-none transition-colors duration-200 truncate ${todo.completed ? 'text-gray-500 line-through' : 'text-gray-200'
                        }`}
                    >
                      {todo.item}
                    </span>
                  </div>

                  <button
                    onClick={() => deleteTodo(todo.id)}
                    className="opacity-0 group-hover:opacity-100 p-2 text-gray-500 hover:text-red-400 transition-all duration-200 focus:opacity-100"
                    aria-label="Delete todo"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Footer */}
        <div className="p-3 bg-gray-900 border-t border-gray-800 text-center">
          <p className="text-xs text-gray-500">{todos.filter(t => !t.completed).length} items left</p>
        </div>
      </div>
    </div>
  )
}

export default App