import { useEffect, useState, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { listProducts, type Product } from '@/lib/api';
import { addToCart, dispatchCartEvent } from '@/lib/cart';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import NavBar from '@/components/NavBar';

const ALL_CATEGORIES = 'All';

export default function Catalog() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const [quantities, setQuantities] = useState<Record<string, number>>({});

  const category = searchParams.get('category') || ALL_CATEGORIES;

  useEffect(() => {
    listProducts()
      .then((data) => {
        const prods = data.products || [];
        setProducts(prods);
        const q: Record<string, number> = {};
        prods.forEach((p) => { q[p.id] = 1; });
        setQuantities(q);
      })
      .catch((err: Error) => setError(err.message));
  }, []);

  const categories = useMemo(() => {
    return [...new Set(products.map((p) => p.category).filter((c): c is string => !!c))];
  }, [products]);

  const filtered = useMemo(() => {
    let result = products;
    if (category !== ALL_CATEGORIES) {
      result = result.filter((p) => p.category === category);
    }
    if (search.trim()) {
      const q = search.toLowerCase();
      result = result.filter((p) =>
        p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q)
      );
    }
    return result;
  }, [products, category, search]);

  function handleAddToCart(product: Product) {
    const qty = quantities[product.id] || 1;
    addToCart({ productId: product.id, name: product.name, price: product.price, quantity: qty });
    dispatchCartEvent();
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background">
        <NavBar />
        <div className="flex items-center justify-center min-h-[60vh]">
          <p className="text-red-500">Failed to load products: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <NavBar showSearch search={search} onSearchChange={setSearch} />

      <div className="max-w-7xl mx-auto px-4 py-8 flex gap-8">
        <aside className="w-48 shrink-0 hidden md:block">
          <nav className="space-y-1 sticky top-24">
            <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3">
              Categories
            </h3>
            <button
              onClick={() => setSearchParams({})}
              className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                category === ALL_CATEGORIES
                    ? 'bg-primary/10 text-primary font-medium'
                    : 'text-muted-foreground hover:bg-primary/5 hover:text-primary'
              }`}
            >
              All Products
            </button>
            {categories.map((cat) => (
              <button
                key={cat}
                onClick={() => setSearchParams({ category: cat })}
                className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                  category === cat
                    ? 'bg-primary/10 text-primary font-medium'
                    : 'text-muted-foreground hover:bg-primary/5 hover:text-primary'
                }`}
              >
                {cat}
              </button>
            ))}
          </nav>
        </aside>

        <div className="flex-1">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-2xl font-bold">
              {category === ALL_CATEGORIES ? 'All Products' : category}
            </h1>
            <p className="text-sm text-muted-foreground">{filtered.length} items</p>
          </div>

          {filtered.length === 0 ? (
            <div className="text-center py-16">
              <p className="text-muted-foreground">No products found</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {filtered.map((p) => (
                <Card key={p.id} className="flex flex-col hover:shadow-lg hover:border-primary/30 transition-all">
                  <CardHeader className="pb-0">
                    <div className="w-full h-40 bg-gradient-to-br from-primary/5 to-accent/10 rounded-lg flex items-center justify-center">
                      <span className="text-4xl text-muted-foreground/30">
                        {p.name.charAt(0)}
                      </span>
                    </div>
                  </CardHeader>
                  <CardContent className="flex-1 flex flex-col">
                    <CardTitle className="text-base mt-3">{p.name}</CardTitle>
                    <p className="text-xs text-muted-foreground line-clamp-2 mt-1">{p.description}</p>
                    <div className="flex items-center justify-between mt-auto pt-3">
                      <span className="text-xl font-bold">${p.price?.toFixed(2)}</span>
                      <Badge variant={p.availableStock > 0 ? 'outline' : 'destructive'}>
                        {p.availableStock > 0 ? `${p.availableStock} in stock` : 'Out of stock'}
                      </Badge>
                    </div>
                  </CardContent>
                  <CardFooter className="pt-0 gap-2">
                    <Input
                      type="number"
                      min="1"
                      max={p.availableStock}
                      value={quantities[p.id] || 1}
                      onChange={(e) => setQuantities((prev) => ({ ...prev, [p.id]: Math.max(1, Number(e.target.value)) }))}
                      className="w-16 h-9 text-sm"
                    />
                    <Button
                      className="flex-1"
                      size="sm"
                      disabled={p.availableStock <= 0}
                      onClick={() => handleAddToCart(p)}
                    >
                      Add to Cart
                    </Button>
                  </CardFooter>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
