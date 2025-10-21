import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { MessageSquare, Plus, Users, Clock, Trash2, CheckCircle, ArrowRight, RotateCcw } from 'lucide-react';
import { retrospectivesAPI } from '../services/api';
import { useAuth } from '../services/AuthContext';
import toast from 'react-hot-toast';
import ConfirmModal from '../components/ConfirmModal';

const RetrospectivesPage = () => {
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const { data: retrospectives, isLoading: retrospectivesLoading } = useQuery('userRetrospectives', retrospectivesAPI.getRetrospectives);
  
  // Modal states
  const [showEndModal, setShowEndModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showReopenModal, setShowReopenModal] = useState(false);
  const [selectedRetrospective, setSelectedRetrospective] = useState(null);

  // Check if current user is the creator of the retrospective
  const isRetrospectiveOwner = (retrospective) => {
    return user && retrospective.created_by === user.id;
  };

  const endRetrospectiveMutation = useMutation(
    (id) => retrospectivesAPI.endRetrospective(id),
    {
      onSuccess: () => {
        queryClient.invalidateQueries('userRetrospectives');
        toast.success('Retrospectiva encerrada!');
      },
      onError: (error) => {
        toast.error('Erro ao encerrar retrospectiva: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const deleteRetrospectiveMutation = useMutation(
    (id) => retrospectivesAPI.deleteRetrospective(id),
    {
      onSuccess: () => {
        queryClient.invalidateQueries('userRetrospectives');
        toast.success('Retrospectiva deletada com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao deletar retrospectiva: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const reopenRetrospectiveMutation = useMutation(
    (id) => retrospectivesAPI.reopenRetrospective(id),
    {
      onSuccess: () => {
        queryClient.invalidateQueries('userRetrospectives');
        toast.success('Retrospectiva reaberta com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao reabrir retrospectiva: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const handleEndRetrospective = (retro) => {
    setSelectedRetrospective(retro);
    setShowEndModal(true);
  };

  const handleDeleteRetrospective = (retro) => {
    setSelectedRetrospective(retro);
    setShowDeleteModal(true);
  };

  const handleReopenRetrospective = (retro) => {
    setSelectedRetrospective(retro);
    setShowReopenModal(true);
  };

  const confirmEndRetrospective = () => {
    if (selectedRetrospective) {
      endRetrospectiveMutation.mutate(selectedRetrospective.id);
      setShowEndModal(false);
      setSelectedRetrospective(null);
    }
  };

  const confirmDeleteRetrospective = () => {
    if (selectedRetrospective) {
      deleteRetrospectiveMutation.mutate(selectedRetrospective.id);
      setShowDeleteModal(false);
      setSelectedRetrospective(null);
    }
  };

  const confirmReopenRetrospective = () => {
    if (selectedRetrospective) {
      reopenRetrospectiveMutation.mutate(selectedRetrospective.id);
      setShowReopenModal(false);
      setSelectedRetrospective(null);
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-800';
      case 'closed': return 'bg-blue-100 text-blue-800';
      case 'planned': return 'bg-yellow-100 text-yellow-800';
      case 'archived': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'active': return 'Em andamento';
      case 'closed': return 'Encerrada';
      case 'planned': return 'Planejada';
      case 'archived': return 'Arquivada';
      default: return status;
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'active': return '🚀';
      case 'closed': return '✅';
      case 'planned': return '📋';
      case 'archived': return '📁';
      default: return '📄';
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-medium text-gray-900">Retrospectivas</h1>
          <p className="text-gray-500 mt-1">Gerencie suas retrospectivas</p>
        </div>
        <Link 
          to="/retrospectives/new" 
          className="btn btn-primary"
        >
          <Plus className="h-4 w-4 mr-2" />
          Nova Retrospectiva
        </Link>
      </div>

      {retrospectivesLoading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-500">Carregando retrospectivas...</p>
        </div>
      ) : retrospectives?.data && retrospectives.data.length > 0 ? (
        <div className="space-y-3">
          {retrospectives.data.map((retro) => (
            <div key={retro.id} className="card">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <div>
                    <h3 className="font-medium text-gray-900">{retro.title}</h3>
                    <div className="flex items-center space-x-2 mt-1">
                      <span className={`badge ${getStatusColor(retro.status)}`}>
                        {getStatusText(retro.status)}
                      </span>
                    </div>
                    <div className="flex items-center space-x-4 text-sm text-gray-500 mt-1">
                      <div className="flex items-center space-x-1">
                        <span className="text-gray-400">Template:</span>
                        <span className="capitalize">{retro.template.replace('_', ' ')}</span>
                      </div>
                      <div className="flex items-center space-x-1">
                        <span className="text-gray-400">Criada em:</span>
                        <span>{new Date(retro.created_at).toLocaleDateString('pt-BR')}</span>
                      </div>
                      <div className="flex items-center space-x-1">
                        <span className="text-gray-400">Participantes:</span>
                        <span>{retro.participants?.length || 0}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                        {retro.status === 'active' && isRetrospectiveOwner(retro) && (
                          <button
                            onClick={() => handleEndRetrospective(retro)}
                            disabled={endRetrospectiveMutation.isLoading}
                            className="flex flex-col items-center space-y-1 text-green-600 hover:text-green-800 p-2 rounded-md hover:bg-green-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            title="Encerrar retrospectiva"
                          >
                            <CheckCircle className="h-4 w-4" />
                            <span className="text-xs text-gray-400">Encerrar</span>
                          </button>
                        )}
                        {retro.status === 'closed' && isRetrospectiveOwner(retro) && (
                          <button
                            onClick={() => handleReopenRetrospective(retro)}
                            disabled={reopenRetrospectiveMutation.isLoading}
                            className="flex flex-col items-center space-y-1 text-blue-600 hover:text-blue-800 p-2 rounded-md hover:bg-blue-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            title="Reabrir retrospectiva"
                          >
                            <RotateCcw className="h-4 w-4" />
                            <span className="text-xs text-gray-400">Reabrir</span>
                          </button>
                        )}
                        {retro.status !== 'closed' && isRetrospectiveOwner(retro) && (
                          <button
                            onClick={() => handleDeleteRetrospective(retro)}
                            disabled={deleteRetrospectiveMutation.isLoading}
                            className="flex flex-col items-center space-y-1 text-red-600 hover:text-red-800 p-2 rounded-md hover:bg-red-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            title="Deletar retrospectiva"
                          >
                            <Trash2 className="h-4 w-4" />
                            <span className="text-xs text-gray-400">Deletar</span>
                          </button>
                        )}
                  <Link
                    to={`/retrospectives/${retro.id}`}
                    className="flex flex-col items-center space-y-1 text-gray-600 hover:text-gray-900 p-2 rounded-md hover:bg-gray-50 transition-colors"
                    title="Entrar na retrospectiva"
                  >
                    <ArrowRight className="h-4 w-4" />
                    <span className="text-xs text-gray-400">Entrar</span>
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <MessageSquare className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            Nenhuma retrospectiva encontrada
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            Comece criando sua primeira retrospectiva.
          </p>
        </div>
      )}

      {/* Confirmation Modals */}
      <ConfirmModal
        isOpen={showEndModal}
        onClose={() => {
          setShowEndModal(false);
          setSelectedRetrospective(null);
        }}
        onConfirm={confirmEndRetrospective}
        title="Encerrar retrospectiva"
        message={`Tem certeza que deseja encerrar a retrospectiva "${selectedRetrospective?.title}"?`}
        confirmText="Encerrar"
        cancelText="Cancelar"
      />

      <ConfirmModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setSelectedRetrospective(null);
        }}
        onConfirm={confirmDeleteRetrospective}
        title="Deletar retrospectiva"
        message={`Tem certeza que deseja deletar a retrospectiva "${selectedRetrospective?.title}"? Esta ação não pode ser desfeita.`}
        confirmText="Deletar"
        cancelText="Cancelar"
      />

      <ConfirmModal
        isOpen={showReopenModal}
        onClose={() => {
          setShowReopenModal(false);
          setSelectedRetrospective(null);
        }}
        onConfirm={confirmReopenRetrospective}
        title="Reabrir retrospectiva"
        message={`Tem certeza que deseja reabrir a retrospectiva "${selectedRetrospective?.title}"?`}
        confirmText="Reabrir"
        cancelText="Cancelar"
      />
    </div>
  );
};

export default RetrospectivesPage;
