import { useState } from 'react'
import './index.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="flex h-screen w-full flex-col items-center justify-center bg-gray-900 text-white">
      <h1 className="mb-4 text-4xl font-bold text-emerald-400">🧙‍♂️ Mestre Stack Vite</h1>
      <p className="mb-8 text-gray-400">React + TypeScript + TailwindCSS</p>
      
      <button 
        onClick={() => setCount((count) => count + 1)}
        className="rounded bg-emerald-600 px-4 py-2 font-bold hover:bg-emerald-500 transition"
      >
        Contador: {count}
      </button>
    </div>
  )
}

export default App
