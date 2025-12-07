import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { signin } from '../api/auth'
import { useAuth } from '../context/AuthContext'

export default function Signin() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { setAuth } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      const data = await signin(email, password)
      setAuth(data.token, data.email)
      navigate('/me')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Signin failed')
    }
  }

  return (
    <div className="container">
      <div className="card">
        <h1>Sign In</h1>
        <form onSubmit={handleSubmit}>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && <div className="error">{error}</div>}
          <button type="submit">Sign In</button>
        </form>
        <p>
          Don't have an account? <Link to="/signup">Sign Up</Link>
        </p>
      </div>
    </div>
  )
}
