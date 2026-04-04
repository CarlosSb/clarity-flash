const USER_ID_KEY = 'user_id'

export function getUserId(): string | null {
  return localStorage.getItem(USER_ID_KEY)
}

export function setUserId(id: string): void {
  localStorage.setItem(USER_ID_KEY, id)
}

export function clearUserId(): void {
  localStorage.removeItem(USER_ID_KEY)
}

export function isAuthenticated(): boolean {
  return !!getUserId()
}

/**
 * Authenticate a user by storing their ID.
 * In production this would involve real login flow.
 * For MVP, we use a simple user_id stored in localStorage.
 */
export function authenticate(userId: string): boolean {
  if (!userId.trim()) return false

  // Basic UUID validation
  const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
  if (!uuidPattern.test(userId)) return false

  setUserId(userId)
  return true
}
