import React from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { Plus, Users, Settings, Trash2, MoreVertical } from 'lucide-react';
import toast from 'react-hot-toast';
import { teamsAPI } from '../services/api';

const TeamsPage = () => {
  const queryClient = useQueryClient();

  const { data: teams, isLoading } = useQuery('userTeams', teamsAPI.getTeams);


  const deleteTeamMutation = useMutation(teamsAPI.deleteTeam, {
    onSuccess: () => {
      queryClient.invalidateQueries('userTeams');
      toast.success('Time removido com sucesso!');
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'Erro ao remover time');
    },
  });


  const handleDeleteTeam = (teamId) => {
    if (window.confirm('Tem certeza que deseja remover este time?')) {
      deleteTeamMutation.mutate(teamId);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Meus Times
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Gerencie seus times e colaboradores
          </p>
        </div>
        <div className="mt-4 flex md:mt-0 md:ml-4">
          <Link
            to="/teams/new"
            className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-purple-400 hover:bg-purple-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-300"
          >
            <Plus className="h-4 w-4 mr-2" />
            Novo Time
          </Link>
        </div>
      </div>

      {/* Teams Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-gray-200 rounded w-full mb-2"></div>
              <div className="h-3 bg-gray-200 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      ) : teams?.data?.length > 0 ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {teams.data.map((team) => (
            <div key={team.id} className="card hover:shadow-md transition-shadow duration-200">
              <div className="flex items-start justify-between">
                <div className="flex items-center space-x-3">
                  <div className="flex-shrink-0">
                    <div className="h-10 w-10 bg-primary-100 rounded-lg flex items-center justify-center">
                      <Users className="h-5 w-5 text-primary-600" />
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-lg font-medium text-gray-900 truncate">
                      {team.name}
                    </h3>
                    <p className="text-sm text-gray-500 truncate">
                      {team.description || 'Sem descrição'}
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  <Link
                    to={`/teams/${team.id}`}
                    className="text-primary-600 hover:text-primary-900"
                  >
                    <Settings className="h-4 w-4" />
                  </Link>
                  <button
                    onClick={() => handleDeleteTeam(team.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
              
              <div className="mt-4 flex items-center justify-between">
                <div className="flex items-center space-x-4 text-sm text-gray-500">
                  <span>Membros: {team.member_count || 0}</span>
                  <span>Retrospectivas: {team.retrospective_count || 0}</span>
                </div>
                <Link
                  to={`/teams/${team.id}`}
                  className="text-primary-600 hover:text-primary-900 text-sm font-medium"
                >
                  Ver detalhes →
                </Link>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <Users className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            Nenhum time encontrado
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            Comece criando seu primeiro time para organizar retrospectivas.
          </p>
          <div className="mt-6">
            <Link
              to="/teams/new"
              className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-purple-400 hover:bg-purple-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-300"
            >
              <Plus className="h-4 w-4 mr-2" />
              Criar Primeiro Time
            </Link>
          </div>
        </div>
      )}

    </div>
  );
};


export default TeamsPage;
