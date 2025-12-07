import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import Signup from './pages/Signup'
import Signin from './pages/Signin'
import Me from './pages/Me'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { token } = useAuth()
  return token ? <>{children}</> : <Navigate to="/signin" />
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/signup" element={<Signup />} />
          <Route path="/signin" element={<Signin />} />
          <Route
            path="/me"
            element={
              <PrivateRoute>
                <Me />
              </PrivateRoute>
            }
          />
          <Route path="/" element={<Navigate to="/signin" />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
