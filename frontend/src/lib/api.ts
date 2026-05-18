export interface Product {
  id: string
  name: string
  description: string
  price: number
  availableStock: number
  category?: string
}

interface ProductsResponse {
  products: Product[]
}

interface AuthResponse {
  token: string
}

interface OrderItemRequest {
  product_id: string
  quantity: number
}

interface OrderItemResponse {
  productId: string
  quantity: number
  price: number
}

export interface OrderResponse {
  id: string
  userId: string
  status: string
  orderDate?: { seconds: number }
  items: OrderItemResponse[]
}

const BASE = '/api'

function getToken(): string | null {
  return localStorage.getItem('token')
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = { 'Content-Type': 'application/json', ...(options.headers as Record<string, string>) }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE}${path}`, { ...options, headers })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || res.statusText)
  }
  return res.json()
}

export function login(userId: string, password: string): Promise<AuthResponse> {
  return request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, password }),
  })
}

export function register(userId: string, password: string): Promise<AuthResponse> {
  return request('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, password }),
  })
}

export function listProducts(category?: string): Promise<ProductsResponse> {
  const params = category ? `?category=${encodeURIComponent(category)}` : ''
  return request(`/products${params}`)
}

export function createOrder(userId: string, items: OrderItemRequest[], shippingAddress: string): Promise<OrderResponse> {
  return request('/orders', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, items, shipping_address: shippingAddress }),
  })
}

export function getOrder(id: string): Promise<OrderResponse> {
  return request(`/orders/${id}`)
}
