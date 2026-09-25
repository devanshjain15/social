import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router'
import App from './App.tsx'
import './index.css'
import { ConfirmUser } from './pages/ConfirmUser.tsx'

export const BASE_API = "http://localhost:8000/v1"

const router = createBrowserRouter([
  {
    path: "/", 
    element: <App />
  }, {
    path: "/confirm/:token", 
    element: <ConfirmUser/>
  }
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
