import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { listProducts, createOrder } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

export default function CreateOrder() {
  const navigate = useNavigate();
  const [products, setProducts] = useState([]);
  const [userId, setUserId] = useState('');
  const [shippingAddress, setShippingAddress] = useState('');
  const [items, setItems] = useState([{ productId: '', quantity: 1 }]);
  const [error, setError] = useState('');

  useEffect(() => {
    listProducts()
      .then((data) => setProducts(data.products || []))
      .catch((err) => setError(err.message));
  }, []);

  function addItem() {
    setItems([...items, { productId: '', quantity: 1 }]);
  }

  function updateItem(index, field, value) {
    const next = [...items];
    next[index][field] = value;
    setItems(next);
  }

  function removeItem(index) {
    setItems(items.filter((_, i) => i !== index));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    try {
      const orderItems = items.map((item) => ({
        product_id: item.productId,
        quantity: Number(item.quantity),
      }));
      const order = await createOrder(userId, orderItems, shippingAddress);
      navigate(`/orders/${order.id}`);
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <Card>
        <CardHeader>
          <CardTitle>Create Order</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              placeholder="User ID"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              required
            />
            <Input
              placeholder="Shipping Address"
              value={shippingAddress}
              onChange={(e) => setShippingAddress(e.target.value)}
              required
            />
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Items</span>
                <Button type="button" variant="outline" size="sm" onClick={addItem}>
                  Add Item
                </Button>
              </div>
              {items.map((item, i) => (
                <div key={i} className="flex gap-2 items-end">
                  <Select
                    value={item.productId}
                    onValueChange={(v) => updateItem(i, 'productId', v)}
                  >
                    <SelectTrigger className="flex-1">
                      <SelectValue placeholder="Select product" />
                    </SelectTrigger>
                    <SelectContent>
                      {products.map((p) => (
                        <SelectItem key={p.id} value={p.id}>
                          {p.name} (${p.price?.toFixed(2)})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Input
                    type="number"
                    min="1"
                    className="w-20"
                    value={item.quantity}
                    onChange={(e) => updateItem(i, 'quantity', e.target.value)}
                    required
                  />
                  {items.length > 1 && (
                    <Button type="button" variant="ghost" size="sm" onClick={() => removeItem(i)}>
                      X
                    </Button>
                  )}
                </div>
              ))}
            </div>
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button type="submit" className="w-full">Place Order</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
