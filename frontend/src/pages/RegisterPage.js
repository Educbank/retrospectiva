import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { MessageSquare, Mail, Lock, User } from 'lucide-react';
import { useAuth } from '../services/AuthContext';
import InputField from '../components/InputField';
import LoadingSpinner from '../components/LoadingSpinner';

const RegisterPage = () => {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: '',
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Basic validation
    if (formData.password !== formData.confirmPassword) {
      setErrors({ confirmPassword: 'As senhas não coincidem' });
      return;
    }
    
    setLoading(true);

    try {
      const { confirmPassword, ...registerData } = formData;
      const result = await register(registerData);
      if (result.success) {
        navigate('/dashboard');
      }
    } catch (error) {
      console.error('Register error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <div className="mx-auto h-12 w-12 flex items-center justify-center bg-primary-600 rounded-lg">
            <MessageSquare className="h-8 w-8 text-white" />
          </div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Crie sua conta
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Ou{' '}
            <Link
              to="/login"
              className="font-medium text-primary-600 hover:text-primary-500"
            >
              entre na sua conta existente
            </Link>
          </p>
        </div>
        
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <InputField
              label="Nome completo"
              name="name"
              type="text"
              placeholder="Seu nome completo"
              value={formData.name}
              onChange={handleChange}
              error={errors.name}
              required
              icon={User}
              autoComplete="name"
            />

            <InputField
              label="Email"
              name="email"
              type="email"
              placeholder="seu@email.com"
              value={formData.email}
              onChange={handleChange}
              error={errors.email}
              required
              icon={Mail}
              autoComplete="email"
            />
            
            <InputField
              label="Senha"
              name="password"
              type="password"
              placeholder="Mínimo 6 caracteres"
              value={formData.password}
              onChange={handleChange}
              error={errors.password}
              required
              icon={Lock}
              autoComplete="new-password"
            />

            <InputField
              label="Confirmar senha"
              name="confirmPassword"
              type="password"
              placeholder="Confirme sua senha"
              value={formData.confirmPassword}
              onChange={handleChange}
              error={errors.confirmPassword}
              required
              icon={Lock}
              autoComplete="new-password"
            />
          </div>

          <div>
            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary btn-lg w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <LoadingSpinner size="sm" color="white" text="Criando conta..." />
              ) : (
                'Criar conta'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;
