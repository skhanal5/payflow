export interface CartItem {
  productId: string
  name: string
  price: number
  quantity: number
}

const CART_KEY = 'cart'

export function getCart(): CartItem[] {
  try {
    return JSON.parse(localStorage.getItem(CART_KEY) || '[]')
  } catch {
    return []
  }
}

export function addToCart(item: CartItem): void {
  const cart = getCart()
  const existing = cart.find((i) => i.productId === item.productId)
  if (existing) {
    existing.quantity += item.quantity
  } else {
    cart.push(item)
  }
  localStorage.setItem(CART_KEY, JSON.stringify(cart))
}

export function updateQuantity(productId: string, quantity: number): void {
  const cart = getCart()
  const item = cart.find((i) => i.productId === productId)
  if (item) {
    item.quantity = quantity
  }
  localStorage.setItem(CART_KEY, JSON.stringify(cart))
}

export function removeFromCart(productId: string): void {
  const cart = getCart().filter((i) => i.productId !== productId)
  localStorage.setItem(CART_KEY, JSON.stringify(cart))
}

export function clearCart(): void {
  localStorage.setItem(CART_KEY, '[]')
}

export function cartCount(): number {
  return getCart().reduce((sum, i) => sum + i.quantity, 0)
}

export function cartTotal(): number {
  return getCart().reduce((sum, i) => sum + i.price * i.quantity, 0)
}

export function dispatchCartEvent(): void {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new Event('cart-updated'))
  }
}
