import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listProducts, type Product } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

export default function Catalog() {
  const [products, setProducts] = useState<Product[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    listProducts()
      .then((data) => setProducts(data.products || []))
      .catch((err: Error) => setError(err.message));
  }, []);

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-red-500">Failed to load products: {error}</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Products</h1>
        <Link to="/order/new">
          <Button>New Order</Button>
        </Link>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {products.map((p) => (
          <Card key={p.id}>
            <CardHeader>
              <CardTitle>{p.name}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <p className="text-sm text-muted-foreground">{p.description}</p>
              <div className="flex items-center justify-between">
                <span className="text-lg font-semibold">${p.price?.toFixed(2)}</span>
                <Badge variant={p.availableStock > 0 ? 'default' : 'destructive'}>
                  {p.availableStock > 0 ? `${p.availableStock} in stock` : 'Out of stock'}
                </Badge>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
