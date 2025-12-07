import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { signout } from '../api/auth'

export default function Me() {
  const { email, clearAuth } = useAuth()
  const navigate = useNavigate()

  const handleSignout = async () => {
    await signout()
    clearAuth()
    navigate('/signin')
  }

  return (
    <div className="container">
      <div className="card">
        <h1>Me</h1>
        <p>This is a protected page</p>
        <p>Email: {email}</p>
        <button onClick={handleSignout}>Sign Out</button>
      </div>
    </div>
  )
}
