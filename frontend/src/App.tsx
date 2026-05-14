import { ReactNode } from 'react';
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import Catalog from '@/pages/Catalog';
import CreateOrder from '@/pages/CreateOrder';
import OrderDetail from '@/pages/OrderDetail';
import Login from '@/pages/Login';
import Register from '@/pages/Register';

function ProtectedRoute({ children }: { children: ReactNode }) {
  const token = localStorage.getItem('token');
  if (!token) return <Navigate to="/login" replace />;
  return children;
}

function Layout({ children }: { children: ReactNode }) {
  const token = localStorage.getItem('token');
  return (
    <div>
      {token && (
        <header className="border-b px-6 py-3 flex items-center justify-between">
          <Link to="/" className="font-semibold">Payflow</Link>
          <div className="flex items-center gap-2">
            <Link to="/order/new"><Button variant="ghost" size="sm">New Order</Button></Link>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => { localStorage.removeItem('token'); window.location.href = '/login'; }}
            >
              Logout
            </Button>
          </div>
        </header>
      )}
      <main>{children}</main>
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/" element={<Catalog />} />
          <Route path="/order/new" element={<ProtectedRoute><CreateOrder /></ProtectedRoute>} />
          <Route path="/orders/:id" element={<ProtectedRoute><OrderDetail /></ProtectedRoute>} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}
