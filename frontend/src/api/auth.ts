const API_BASE = '/api'

export async function signup(email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE}/auth/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || 'Signup failed')
  }
}

export async function signin(email: string, password: string): Promise<{ token: string; email: string }> {
  const response = await fetch(`${API_BASE}/auth/signin`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || 'Signin failed')
  }

  return response.json()
}

export async function getMe(token: string): Promise<{ email: string }> {
  const response = await fetch(`${API_BASE}/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })

  if (!response.ok) {
    throw new Error('Failed to fetch user info')
  }

  return response.json()
}

export async function signout(): Promise<void> {
  await fetch(`${API_BASE}/auth/signout`, { method: 'POST' })
}
