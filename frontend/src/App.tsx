import { ReactNode } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Catalog from '@/pages/Catalog';
import Cart from '@/pages/Cart';
import CreateOrder from '@/pages/CreateOrder';
import OrderDetail from '@/pages/OrderDetail';
import Login from '@/pages/Login';
import Register from '@/pages/Register';

function ProtectedRoute({ children }: { children: ReactNode }) {
  const token = localStorage.getItem('token');
  if (!token) return <Navigate to="/login" replace />;
  return children;
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<Catalog />} />
        <Route path="/cart" element={<Cart />} />
        <Route path="/order/new" element={<ProtectedRoute><CreateOrder /></ProtectedRoute>} />
        <Route path="/orders/:id" element={<ProtectedRoute><OrderDetail /></ProtectedRoute>} />
      </Routes>
    </BrowserRouter>
  );
}
