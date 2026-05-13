import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { register } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function Register() {
  const navigate = useNavigate();
  const [userId, setUserId] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  async function handleSubmit(e) {
    e.preventDefault();
    try {
      const data = await register(userId, password);
      localStorage.setItem('token', data.token);
      navigate('/');
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <Card className="w-96">
        <CardHeader>
          <CardTitle>Register</CardTitle>
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
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button type="submit" className="w-full">Register</Button>
            <p className="text-sm text-center">
              Already have an account? <Link to="/login" className="underline">Login</Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
