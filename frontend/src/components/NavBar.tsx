import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cartCount, dispatchCartEvent } from '@/lib/cart';

interface NavBarProps {
  search?: string
  onSearchChange?: (value: string) => void
  showSearch?: boolean
}

export default function NavBar({ search, onSearchChange, showSearch }: NavBarProps) {
  const navigate = useNavigate();
  const [count, setCount] = useState(cartCount());
  const token = localStorage.getItem('token');

  useEffect(() => {
    const handler = () => setCount(cartCount());
    window.addEventListener('storage', handler);
    window.addEventListener('cart-updated', handler);
    return () => {
      window.removeEventListener('storage', handler);
      window.removeEventListener('cart-updated', handler);
    };
  }, []);

  return (
    <header className="sticky top-0 z-50 border-b bg-white/80 backdrop-blur-md">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center gap-4">
        <Link to="/" className="text-xl font-bold tracking-tight text-primary shrink-0">
          Payflow
        </Link>

        {showSearch && (
          <div className="flex-1 max-w-md mx-auto">
            <Input
              placeholder="Search products..."
              value={search || ''}
              onChange={(e) => onSearchChange?.(e.target.value)}
              className="bg-muted/50"
            />
          </div>
        )}

        <div className="flex items-center gap-2 ml-auto shrink-0">
          <Link to="/cart">
            <Button variant="ghost" size="sm" className="relative gap-1">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="8" cy="21" r="1"/><circle cx="21" cy="21" r="1"/>
                <path d="M3 3h2l.4 2M7 13h10l4-8H5.4"/><path d="M7 13 5.4 5H3"/>
              </svg>
              Cart
              {count > 0 && (
                <span className="absolute -top-1 -right-1 bg-primary text-primary-foreground text-[10px] font-bold rounded-full min-w-[18px] h-[18px] flex items-center justify-center">
                  {count}
                </span>
              )}
            </Button>
          </Link>

          {token ? (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                localStorage.removeItem('token');
                dispatchCartEvent();
                navigate('/');
              }}
            >
              Sign Out
            </Button>
          ) : (
            <Link to="/login">
              <Button variant="default" size="sm">Sign In</Button>
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}


