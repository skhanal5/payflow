import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getCart, updateQuantity, removeFromCart, clearCart, cartTotal, dispatchCartEvent, type CartItem } from '@/lib/cart';
import { createOrder } from '@/lib/api';
import NavBar from '@/components/NavBar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function Cart() {
  const navigate = useNavigate();
  const [items, setItems] = useState<CartItem[]>([]);
  const [shippingAddress, setShippingAddress] = useState('');
  const [error, setError] = useState('');
  const [placing, setPlacing] = useState(false);
  const token = localStorage.getItem('token');

  useEffect(() => {
    setItems(getCart());
  }, []);

  function refresh() {
    setItems(getCart());
    dispatchCartEvent();
  }

  function handleQtyChange(productId: string, qty: number) {
    if (qty <= 0) {
      removeFromCart(productId);
    } else {
      updateQuantity(productId, qty);
    }
    refresh();
  }

  function handleRemove(productId: string) {
    removeFromCart(productId);
    refresh();
  }

  async function handlePlaceOrder() {
    if (!token) {
      navigate('/login');
      return;
    }
    if (!shippingAddress.trim()) {
      setError('Please enter a shipping address');
      return;
    }
    setPlacing(true);
    setError('');
    try {
      const userId = ''; // will be extracted from token by backend
      const orderItems = items.map((i) => ({
        product_id: i.productId,
        quantity: i.quantity,
      }));
      const order = await createOrder(userId, orderItems, shippingAddress);
      clearCart();
      refresh();
      navigate(`/orders/${order.id}`);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setPlacing(false);
    }
  }

  const total = cartTotal();

  return (
    <div className="min-h-screen bg-background">
      <NavBar />

      <div className="max-w-3xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6">Shopping Cart</h1>

        {items.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-muted-foreground mb-4">Your cart is empty</p>
            <Button onClick={() => navigate('/')}>Continue Shopping</Button>
          </div>
        ) : (
          <div className="space-y-4">
            {items.map((item) => (
              <Card key={item.productId}>
                <CardContent className="flex items-center gap-4 py-4">
                  <div className="w-16 h-16 bg-gradient-to-br from-primary/5 to-accent/10 rounded-lg flex items-center justify-center shrink-0">
                    <span className="text-xl text-muted-foreground/30">{item.name.charAt(0)}</span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{item.name}</p>
                    <p className="text-sm text-muted-foreground">${item.price.toFixed(2)} each</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Input
                      type="number"
                      min="0"
                      value={item.quantity}
                      onChange={(e) => handleQtyChange(item.productId, Math.max(0, Number(e.target.value)))}
                      className="w-16 h-9 text-sm"
                    />
                    <p className="text-sm font-medium w-20 text-right">
                      ${(item.price * item.quantity).toFixed(2)}
                    </p>
                    <Button variant="ghost" size="sm" onClick={() => handleRemove(item.productId)}>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                      </svg>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}

            <Card className="mt-6">
              <CardHeader>
                <CardTitle>Checkout</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {!token && (
                  <p className="text-sm text-amber-600 bg-amber-50 border border-amber-200 rounded-md px-3 py-2">
                    You need to sign in to place an order.
                  </p>
                )}
                <Input
                  placeholder="Shipping Address"
                  value={shippingAddress}
                  onChange={(e) => setShippingAddress(e.target.value)}
                  required
                />
                <div className="flex items-center justify-between text-lg font-bold">
                  <span>Total</span>
                  <span>${total.toFixed(2)}</span>
                </div>
                {error && <p className="text-sm text-red-500">{error}</p>}
                <Button
                  className="w-full"
                  size="lg"
                  disabled={placing || items.length === 0}
                  onClick={handlePlaceOrder}
                >
                  {placing ? 'Placing Order...' : token ? 'Place Order' : 'Sign In to Place Order'}
                </Button>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
