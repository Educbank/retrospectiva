import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { User, Mail, Save, Camera } from 'lucide-react';
import toast from 'react-hot-toast';
import { usersAPI } from '../services/api';

const ProfilePage = () => {
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    avatar: '',
  });

  const queryClient = useQueryClient();
  const { data: user, isLoading } = useQuery('userProfile', usersAPI.getProfile);

  const updateProfileMutation = useMutation(usersAPI.updateProfile, {
    onSuccess: () => {
      queryClient.invalidateQueries('userProfile');
      setIsEditing(false);
      toast.success('Perfil atualizado com sucesso!');
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'Erro ao atualizar perfil');
    },
  });

  const handleEdit = () => {
    if (user?.data) {
      setFormData({
        name: user.data.name,
        email: user.data.email,
        avatar: user.data.avatar || '',
      });
      setIsEditing(true);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    updateProfileMutation.mutate({
      name: formData.name,
      avatar: formData.avatar,
    });
  };

  const handleCancel = () => {
    setIsEditing(false);
    setFormData({
      name: '',
      email: '',
      avatar: '',
    });
  };

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-8 bg-gray-200 rounded w-1/4"></div>
        <div className="card">
          <div className="h-32 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (!user?.data) {
    return (
      <div className="text-center py-12">
        <User className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">
          Perfil não encontrado
        </h3>
      </div>
    );
  }

  const userData = user.data;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Meu Perfil
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Gerencie suas informações pessoais
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Profile Info */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Informações Pessoais
              </h3>
            </div>
            
            {!isEditing ? (
              <div className="space-y-6">
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <div className="h-20 w-20 bg-primary-600 rounded-full flex items-center justify-center">
                      {userData.avatar ? (
                        <img
                          src={userData.avatar}
                          alt="Avatar"
                          className="h-20 w-20 rounded-full object-cover"
                        />
                      ) : (
                        <span className="text-2xl font-medium text-white">
                          {userData.name?.charAt(0).toUpperCase()}
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex-1">
                    <h4 className="text-xl font-medium text-gray-900">
                      {userData.name}
                    </h4>
                    <p className="text-sm text-gray-500">{userData.email}</p>
                    <p className="text-sm text-gray-500">
                      Membro desde {new Date(userData.created_at).toLocaleDateString('pt-BR')}
                    </p>
                  </div>
                </div>
                
                <div className="pt-6 border-t border-gray-200">
                  <button
                    onClick={handleEdit}
                    className="btn btn-primary"
                  >
                    Editar Perfil
                  </button>
                </div>
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <div className="h-20 w-20 bg-primary-600 rounded-full flex items-center justify-center relative">
                      {formData.avatar ? (
                        <img
                          src={formData.avatar}
                          alt="Avatar"
                          className="h-20 w-20 rounded-full object-cover"
                        />
                      ) : (
                        <span className="text-2xl font-medium text-white">
                          {formData.name?.charAt(0).toUpperCase() || 'U'}
                        </span>
                      )}
                      <button
                        type="button"
                        className="absolute inset-0 bg-black bg-opacity-50 rounded-full flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity"
                      >
                        <Camera className="h-6 w-6 text-white" />
                      </button>
                    </div>
                  </div>
                  <div className="flex-1">
                    <div className="space-y-4">
                      <div>
                        <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                          Nome completo
                        </label>
                        <input
                          id="name"
                          name="name"
                          type="text"
                          required
                          value={formData.name}
                          onChange={handleChange}
                          className="input mt-1"
                          placeholder="Seu nome completo"
                        />
                      </div>
                      
                      <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                          Email
                        </label>
                        <input
                          id="email"
                          name="email"
                          type="email"
                          disabled
                          value={formData.email}
                          className="input mt-1 bg-gray-50"
                          placeholder="seu@email.com"
                        />
                        <p className="mt-1 text-xs text-gray-500">
                          O email não pode ser alterado
                        </p>
                      </div>
                      
                      <div>
                        <label htmlFor="avatar" className="block text-sm font-medium text-gray-700">
                          URL do Avatar (opcional)
                        </label>
                        <input
                          id="avatar"
                          name="avatar"
                          type="url"
                          value={formData.avatar}
                          onChange={handleChange}
                          className="input mt-1"
                          placeholder="https://exemplo.com/foto.jpg"
                        />
                      </div>
                    </div>
                  </div>
                </div>
                
                <div className="pt-6 border-t border-gray-200 flex items-center justify-end space-x-3">
                  <button
                    type="button"
                    onClick={handleCancel}
                    className="btn btn-secondary"
                  >
                    Cancelar
                  </button>
                  <button
                    type="submit"
                    disabled={updateProfileMutation.isLoading}
                    className="btn btn-primary"
                  >
                    {updateProfileMutation.isLoading ? (
                      <div className="flex items-center">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                        Salvando...
                      </div>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Salvar
                      </>
                    )}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>

        {/* Account Stats */}
        <div className="space-y-6">
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Estatísticas da Conta
              </h3>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <User className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Times</span>
                </div>
                <span className="text-sm font-medium text-gray-900">3</span>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Mail className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Retrospectivas</span>
                </div>
                <span className="text-sm font-medium text-gray-900">24</span>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <User className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Action Items</span>
                </div>
                <span className="text-sm font-medium text-gray-900">67</span>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Configurações
              </h3>
            </div>
            <div className="space-y-3">
              <button className="w-full btn btn-secondary btn-sm">
                Alterar Senha
              </button>
              <button className="w-full btn btn-secondary btn-sm">
                Notificações
              </button>
              <button className="w-full btn btn-secondary btn-sm">
                Privacidade
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
