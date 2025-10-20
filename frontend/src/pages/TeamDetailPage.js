import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { 
  Users, 
  Plus, 
  Mail, 
  Settings, 
  Trash2, 
  Crown,
  User,
  Calendar,
  MessageSquare,
  TrendingUp
} from 'lucide-react';
import toast from 'react-hot-toast';
import { teamsAPI } from '../services/api';

const TeamDetailPage = () => {
  const { id } = useParams();
  const [showInviteModal, setShowInviteModal] = useState(false);
  const queryClient = useQueryClient();

  const { data: team, isLoading } = useQuery(
    ['team', id],
    () => teamsAPI.getTeam(id),
    { enabled: !!id }
  );

  const { data: analytics } = useQuery(
    ['teamAnalytics', id],
    () => teamsAPI.getAnalytics(id),
    { enabled: !!id }
  );

  const removeMemberMutation = useMutation(
    ({ teamId, userId }) => teamsAPI.removeMember(teamId, userId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['team', id]);
        toast.success('Membro removido com sucesso!');
      },
      onError: (error) => {
        toast.error(error.response?.data?.error || 'Erro ao remover membro');
      },
    }
  );

  const inviteMemberMutation = useMutation(
    ({ teamId, data }) => teamsAPI.addMember(teamId, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['team', id]);
        setShowInviteModal(false);
        toast.success('Convite enviado com sucesso!');
      },
      onError: (error) => {
        toast.error(error.response?.data?.error || 'Erro ao convidar membro');
      },
    }
  );

  const handleRemoveMember = (userId) => {
    if (window.confirm('Tem certeza que deseja remover este membro?')) {
      removeMemberMutation.mutate({ teamId: id, userId });
    }
  };

  const handleInviteMember = (inviteData) => {
    inviteMemberMutation.mutate({ teamId: id, data: inviteData });
  };

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-8 bg-gray-200 rounded w-1/4"></div>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <div className="lg:col-span-2 space-y-6">
            <div className="h-64 bg-gray-200 rounded-lg"></div>
            <div className="h-48 bg-gray-200 rounded-lg"></div>
          </div>
          <div className="h-64 bg-gray-200 rounded-lg"></div>
        </div>
      </div>
    );
  }

  if (!team?.data) {
    return (
      <div className="text-center py-12">
        <Users className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">
          Time não encontrado
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          O time que você está procurando não existe ou você não tem acesso.
        </p>
      </div>
    );
  }

  const teamData = team.data;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            {teamData.name}
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            {teamData.description || 'Sem descrição'}
          </p>
        </div>
        <div className="mt-4 flex space-x-3 md:mt-0 md:ml-4">
          <button
            onClick={() => setShowInviteModal(true)}
            className="btn btn-primary"
          >
            <Plus className="h-4 w-4 mr-2" />
            Convidar Membro
          </button>
          <Link
            to={`/teams/${id}/settings`}
            className="btn btn-secondary"
          >
            <Settings className="h-4 w-4 mr-2" />
            Configurações
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Team Members */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Membros do Time ({teamData.members?.length || 0})
              </h3>
            </div>
            <div className="space-y-4">
              {teamData.members?.map((member) => (
                <div key={member.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className="flex-shrink-0">
                      <div className={`h-10 w-10 rounded-full flex items-center justify-center ${
                        member.role === 'owner' ? 'bg-yellow-100' : 'bg-gray-100'
                      }`}>
                        <User className={`h-5 w-5 ${
                          member.role === 'owner' ? 'text-yellow-600' : 'text-gray-600'
                        }`} />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center space-x-2">
                        <p className="text-sm font-medium text-gray-900">
                          {member.user_name}
                        </p>
                        {member.role === 'owner' && (
                          <Crown className="h-4 w-4 text-yellow-500" />
                        )}
                      </div>
                      <p className="text-sm text-gray-500">
                        {member.role === 'owner' ? 'Proprietário' : 
                         member.role === 'member' ? 'Membro' : 'Visualizador'}
                      </p>
                    </div>
                  </div>
                  {member.role !== 'owner' && (
                    <button
                      onClick={() => handleRemoveMember(member.user_id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Team Stats */}
        <div className="space-y-6">
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Estatísticas
              </h3>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Users className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Membros</span>
                </div>
                <span className="text-sm font-medium text-gray-900">
                  {teamData.members?.length || 0}
                </span>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <MessageSquare className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Retrospectivas</span>
                </div>
                <span className="text-sm font-medium text-gray-900">
                  {analytics?.data?.total_retrospectives || 0}
                </span>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <TrendingUp className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Action Items</span>
                </div>
                <span className="text-sm font-medium text-gray-900">
                  {analytics?.data?.total_action_items || 0}
                </span>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Calendar className="h-5 w-5 text-gray-400" />
                  <span className="text-sm text-gray-600">Criado em</span>
                </div>
                <span className="text-sm font-medium text-gray-900">
                  {teamData.created_at ? new Date(teamData.created_at).toLocaleDateString('pt-BR') : 'N/A'}
                </span>
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">
                Ações Rápidas
              </h3>
            </div>
            <div className="space-y-3">
              <Link
                to={`/retrospectives/new?team=${id}`}
                className="w-full btn btn-primary btn-sm"
              >
                <Plus className="h-4 w-4 mr-2" />
                Nova Retrospectiva
              </Link>
              <Link
                to={`/retrospectives?team=${id}`}
                className="w-full btn btn-secondary btn-sm"
              >
                <MessageSquare className="h-4 w-4 mr-2" />
                Ver Retrospectivas
              </Link>
              <Link
                to={`/analytics?team=${id}`}
                className="w-full btn btn-secondary btn-sm"
              >
                <TrendingUp className="h-4 w-4 mr-2" />
                Ver Analytics
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Invite Member Modal */}
      {showInviteModal && (
        <InviteMemberModal
          onClose={() => setShowInviteModal(false)}
          onSubmit={handleInviteMember}
          isLoading={inviteMemberMutation.isLoading}
        />
      )}
    </div>
  );
};

const InviteMemberModal = ({ onClose, onSubmit, isLoading }) => {
  const [formData, setFormData] = useState({
    email: '',
    role: 'member',
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium text-gray-900">
              Convidar Membro
            </h3>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
            >
              ×
            </button>
          </div>
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email do membro
              </label>
              <input
                id="email"
                name="email"
                type="email"
                required
                value={formData.email}
                onChange={handleChange}
                className="input mt-1"
                placeholder="membro@exemplo.com"
              />
            </div>
            
            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-700">
                Função
              </label>
              <select
                id="role"
                name="role"
                value={formData.role}
                onChange={handleChange}
                className="input mt-1"
              >
                <option value="member">Membro</option>
                <option value="viewer">Visualizador</option>
              </select>
            </div>
            
            <div className="flex items-center justify-end space-x-3 pt-4">
              <button
                type="button"
                onClick={onClose}
                className="btn btn-secondary"
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={isLoading}
                className="btn btn-primary"
              >
                {isLoading ? 'Enviando...' : 'Convidar'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default TeamDetailPage;
