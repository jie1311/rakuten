import { createContext, useContext, useState, useEffect, ReactNode } from 'react'

interface AuthContextType {
  token: string | null
  email: string | null
  setAuth: (token: string, email: string) => void
  clearAuth: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('token')
  )
  const [email, setEmail] = useState<string | null>(
    localStorage.getItem('email')
  )

  const setAuth = (newToken: string, newEmail: string) => {
    localStorage.setItem('token', newToken)
    localStorage.setItem('email', newEmail)
    setToken(newToken)
    setEmail(newEmail)
  }

  const clearAuth = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('email')
    setToken(null)
    setEmail(null)
  }

  useEffect(() => {
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'token') {
        setToken(e.newValue)
      }
      if (e.key === 'email') {
        setEmail(e.newValue)
      }
    }

    window.addEventListener('storage', handleStorageChange)
    return () => window.removeEventListener('storage', handleStorageChange)
  }, [])

  return (
    <AuthContext.Provider value={{ token, email, setAuth, clearAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
