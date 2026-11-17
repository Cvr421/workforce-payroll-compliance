export class ApiError extends Error {
  constructor(public readonly status: number, public readonly code: string, message: string) {
    super(message)
  }
}

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('complylite_token')
  const headers = new Headers(options.headers)
  if (!(options.body instanceof FormData)) headers.set('Content-Type', 'application/json')
  if (token) headers.set('Authorization', `Bearer ${token}`)
  const response = await fetch(`/api/v1${path}`, { ...options, headers })
  if (response.status === 401) {
    localStorage.removeItem('complylite_token')
    window.dispatchEvent(new CustomEvent('auth:expired'))
  }
  const body = await response.json().catch(() => null) as { error?: { code: string; message: string } } | null
  if (!response.ok) {
    throw new ApiError(response.status, body?.error?.code ?? 'internal', body?.error?.message ?? 'Request failed')
  }
  return body as T
}
