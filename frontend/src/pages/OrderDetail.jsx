import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getOrder } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

export default function OrderDetail() {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    getOrder(id)
      .then(setOrder)
      .catch((err) => setError(err.message));
  }, [id]);

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-red-500">{error}</p>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  const statusVariant = order.status === 'CONFIRMED' ? 'default'
    : order.status === 'FAILED' ? 'destructive'
    : 'secondary';

  return (
    <div className="max-w-2xl mx-auto p-6 space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Order {order.id}</CardTitle>
            <Badge variant={statusVariant}>{order.status}</Badge>
          </div>
        </CardHeader>
        <CardContent className="space-y-2">
          <p className="text-sm text-muted-foreground">User: {order.userId}</p>
          <p className="text-sm text-muted-foreground">
            Date: {order.orderDate ? new Date(order.orderDate.seconds * 1000).toLocaleString() : '-'}
          </p>
          <div className="pt-4">
            <h3 className="text-sm font-medium mb-2">Items</h3>
            <div className="space-y-2">
              {order.items?.map((item, i) => (
                <div key={i} className="flex items-center justify-between border rounded p-2">
                  <span>{item.productId}</span>
                  <span className="text-sm text-muted-foreground">
                    x{item.quantity} @ ${item.price?.toFixed(2)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
