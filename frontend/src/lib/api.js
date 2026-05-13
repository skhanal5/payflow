const BASE = '/api';

function getToken() {
  return localStorage.getItem('token');
}

async function request(path, options = {}) {
  const token = getToken();
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, { ...options, headers });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || res.statusText);
  }
  return res.json();
}

export function login(userId, password) {
  return request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, password }),
  });
}

export function register(userId, password) {
  return request('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, password }),
  });
}

export function listProducts(category) {
  const params = category ? `?category=${encodeURIComponent(category)}` : '';
  return request(`/products${params}`);
}

export function getProduct(id) {
  return request(`/products/${id}`);
}

export function createOrder(userId, items, shippingAddress) {
  return request('/orders', {
    method: 'POST',
    body: JSON.stringify({ user_id: userId, items, shipping_address: shippingAddress }),
  });
}

export function getOrder(id) {
  return request(`/orders/${id}`);
}

export function listOrders(userId) {
  return request(`/orders?user_id=${encodeURIComponent(userId)}`);
}
