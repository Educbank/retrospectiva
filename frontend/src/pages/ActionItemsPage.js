import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { useAuth } from '../services/AuthContext';
import { 
  CheckCircle, 
  Clock, 
  Calendar, 
  Filter, 
  Search, 
  Plus,
  Edit3,
  Trash2,
  MessageSquare,
  AlertCircle,
  X,
  Save,
  User,
  FileText,
  CheckSquare,
  Play,
  Pause
} from 'lucide-react';
import { retrospectivesAPI } from '../services/api';
import toast from 'react-hot-toast';
import ConfirmModal from '../components/ConfirmModal';

const ActionItemsPage = () => {
  const { user } = useAuth();
  const [statusFilter, setStatusFilter] = useState('all');
  const [retrospectiveFilter, setRetrospectiveFilter] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [showCompleteModal, setShowCompleteModal] = useState(false);
  const [showFeedbackModal, setShowFeedbackModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [completingActionItem, setCompletingActionItem] = useState(null);
  const [viewingActionItem, setViewingActionItem] = useState(null);
  const [deletingActionItem, setDeletingActionItem] = useState(null);
  const [completionForm, setCompletionForm] = useState({
    feedback: ''
  });

  const queryClient = useQueryClient();

  // Fetch all retrospectives to get action items
  const { data: retrospectives, isLoading } = useQuery(
    'userRetrospectives',
    retrospectivesAPI.getRetrospectives,
    {
      select: (response) => response.data,
    }
  );

  // Update Action Item mutation
  const updateActionItemMutation = useMutation(
    ({ actionItemId, data }) => retrospectivesAPI.updateActionItem(actionItemId, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries('userRetrospectives');
        queryClient.invalidateQueries(['retrospective']);
        toast.success('Action Item atualizado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao atualizar Action Item: ' + (error.response?.data?.error || error.message));
      }
    }
  );

  // Delete Action Item mutation
  const deleteActionItemMutation = useMutation(
    (actionItemId) => retrospectivesAPI.deleteActionItem(actionItemId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries('userRetrospectives');
        toast.success('Action Item excluído com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao excluir Action Item: ' + (error.response?.data?.error || error.message));
      }
    }
  );

  // Flatten all action items from all retrospectives
  const allActionItems = (retrospectives && Array.isArray(retrospectives)) 
    ? retrospectives.flatMap(retro => 
        retro.action_items?.map(actionItem => ({
          ...actionItem,
          retrospective: {
            id: retro.id,
            title: retro.title,
            template: retro.template,
            status: retro.status
          }
        })) || []
      )
    : [];

  // Filter action items
  const filteredActionItems = allActionItems.filter(actionItem => {
    const matchesStatus = statusFilter === 'all' || actionItem.status === statusFilter;
    const matchesRetrospective = retrospectiveFilter === 'all' || actionItem.retrospective.id === retrospectiveFilter;
    const matchesSearch = searchTerm === '' || 
      actionItem.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      actionItem.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      actionItem.retrospective.title.toLowerCase().includes(searchTerm.toLowerCase());
    
    return matchesStatus && matchesRetrospective && matchesSearch;
  });

  const getStatusColor = (status) => {
    switch (status) {
      case 'done': return 'bg-green-100 text-green-800';
      case 'in_progress': return 'bg-yellow-100 text-yellow-800';
      case 'todo': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'done': return 'Concluído';
      case 'in_progress': return 'Em andamento';
      case 'todo': return 'A fazer';
      default: return status;
    }
  };

  const isOverdue = (dueDate, status) => {
    if (!dueDate || status === 'done') return false;
    const today = new Date();
    // Parse the date string to avoid timezone issues
    const due = new Date(dueDate.split('T')[0] + 'T00:00:00');
    today.setHours(0, 0, 0, 0);
    due.setHours(0, 0, 0, 0);
    return due < today;
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    // Parse the date string to avoid timezone issues
    const date = new Date(dateString.split('T')[0] + 'T00:00:00');
    return date.toLocaleDateString('pt-BR');
  };

  const getDueDateColor = (dueDate, status) => {
    if (!dueDate) return 'text-gray-500';
    if (status === 'done') return 'text-green-600';
    if (isOverdue(dueDate, status)) return 'text-red-600';
    return 'text-blue-600';
  };

  const getDueDateText = (dueDate, status) => {
    if (!dueDate) return '';
    if (status === 'done') return 'Concluído';
    if (isOverdue(dueDate, status)) return 'Atrasado';
    return 'No prazo';
  };

  const getRetrospectiveStatusColor = (status) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-800';
      case 'closed': return 'bg-blue-100 text-blue-800';
      case 'planned': return 'bg-yellow-100 text-yellow-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getRetrospectiveStatusText = (status) => {
    switch (status) {
      case 'active': return 'Ativa';
      case 'closed': return 'Encerrada';
      case 'planned': return 'Planejada';
      default: return status;
    }
  };

  const getTemplateIcon = (template) => {
    switch (template) {
      case 'sailboat': return '⛵';
      case 'start_stop_continue': return '🔄';
      case '4ls': return '📚';
      case 'mad_sad_glad': return '😊';
      default: return '📋';
    }
  };

  // Check if user can edit/delete a specific action item
  const canEditActionItem = (actionItem) => {
    return user && actionItem.created_by === user.id;
  };

  // Handle start action item
  const handleStartActionItem = (actionItemId) => {
    updateActionItemMutation.mutate({
      actionItemId,
      data: { status: 'in_progress' }
    });
  };

  // Handle complete action item
  const handleCompleteActionItem = (actionItem) => {
    setCompletingActionItem(actionItem);
    setCompletionForm({ feedback: '' });
    setShowCompleteModal(true);
  };

  // Handle save completion
  const handleSaveCompletion = (e) => {
    e.preventDefault();
    if (!completingActionItem) return;

    const completionText = `
## ✅ Action Item Concluído - ${new Date().toLocaleDateString('pt-BR')}

**Parecer sobre a conclusão:**
${completionForm.feedback}

---
`;

    const updateData = {
      status: 'done',
      description: (completingActionItem.description || '') + completionText
    };

    updateActionItemMutation.mutate({
      actionItemId: completingActionItem.id,
      data: updateData
    });

    setShowCompleteModal(false);
    setCompletingActionItem(null);
  };

  // Handle view feedback
  const handleViewFeedback = (actionItem) => {
    setViewingActionItem(actionItem);
    setShowFeedbackModal(true);
  };

  // Extract feedback from description
  const extractFeedback = (description) => {
    if (!description) return null;
    
    const feedbackMatch = description.match(/## ✅ Action Item Concluído - .*?\n\n\*\*Parecer sobre a conclusão:\*\*\n(.*?)\n\n---/s);
    return feedbackMatch ? feedbackMatch[1].trim() : null;
  };

  // Extract completion date from description
  const extractCompletionDate = (description) => {
    if (!description) return null;
    
    const dateMatch = description.match(/## ✅ Action Item Concluído - (.*?)\n/);
    return dateMatch ? dateMatch[1].trim() : null;
  };

  // Handle delete action item
  const handleDeleteActionItem = (actionItem) => {
    setDeletingActionItem(actionItem);
    setShowDeleteModal(true);
  };

  const confirmDeleteActionItem = () => {
    if (deletingActionItem) {
      deleteActionItemMutation.mutate(deletingActionItem.id);
      setShowDeleteModal(false);
      setDeletingActionItem(null);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-medium text-gray-900">Action Items</h1>
        <p className="text-gray-500 mt-1">Gerencie todos os seus action items</p>
      </div>

      {/* Filters */}
      <div className="card">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          {/* Search */}
          <div>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input pl-10"
                placeholder="Buscar action items..."
              />
            </div>
          </div>

          {/* Status Filter */}
          <div>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="input"
            >
              <option value="all">Todos os status</option>
              <option value="todo">A fazer</option>
              <option value="in_progress">Em andamento</option>
              <option value="done">Concluídos</option>
            </select>
          </div>

          {/* Retrospective Filter */}
          <div>
            <select
              value={retrospectiveFilter}
              onChange={(e) => setRetrospectiveFilter(e.target.value)}
              className="input"
            >
              <option value="all">Todas as retrospectivas</option>
              {retrospectives && Array.isArray(retrospectives) && retrospectives.map((retro) => (
                <option key={retro.id} value={retro.id}>
                  {retro.title}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="mt-4 text-sm text-gray-500">
          {filteredActionItems.length} de {allActionItems.length} action items
        </div>
      </div>

      {/* Action Items List */}
      {filteredActionItems.length > 0 ? (
        <div className="space-y-3">
          {filteredActionItems.map((actionItem) => (
            <div key={actionItem.id} className="card">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <div>
                    <h3 className="font-medium text-gray-900">{actionItem.title}</h3>
                    <div className="flex items-center space-x-2 mt-1">
                      <span className={`badge ${getStatusColor(actionItem.status)}`}>
                        {getStatusText(actionItem.status)}
                      </span>
                      {isOverdue(actionItem.due_date, actionItem.status) && (
                        <span className="badge badge-danger">
                          <AlertCircle className="h-3 w-3 inline mr-1" />
                          Atrasado
                        </span>
                      )}
                    </div>
                    <div className="flex items-center space-x-4 text-sm text-gray-500 mt-1">
                      <div className="flex items-center space-x-1">
                        <span className="text-gray-400">Retro:</span>
                        <Link
                          to={`/retrospectives/${actionItem.retrospective.id}`}
                          className="text-blue-600 hover:text-blue-800"
                        >
                          {actionItem.retrospective.title}
                        </Link>
                      </div>
                      {actionItem.due_date && (
                        <div className="flex items-center space-x-1">
                          <span className="text-gray-400">Prazo:</span>
                          <span className={getDueDateColor(actionItem.due_date, actionItem.status)}>
                            {formatDate(actionItem.due_date)}
                          </span>
                        </div>
                      )}
                      {actionItem.status === 'done' && extractCompletionDate(actionItem.description) && (
                        <div className="flex items-center space-x-1">
                          <span className="text-green-600">
                            (concluído em: {extractCompletionDate(actionItem.description)})
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  {/* Start Button - only for todo status */}
                  {actionItem.status === 'todo' && canEditActionItem(actionItem) && (
                    <button
                      onClick={() => handleStartActionItem(actionItem.id)}
                      className="flex flex-col items-center space-y-1 text-green-600 hover:text-green-800 p-2 rounded-md hover:bg-green-50 transition-colors"
                      title="Iniciar Action Item"
                    >
                      <Play className="h-4 w-4" />
                      <span className="text-xs text-gray-400">Iniciar</span>
                    </button>
                  )}

                  {/* Complete Button - only for in_progress status */}
                  {actionItem.status === 'in_progress' && canEditActionItem(actionItem) && (
                    <button
                      onClick={() => handleCompleteActionItem(actionItem)}
                      className="flex flex-col items-center space-y-1 text-blue-600 hover:text-blue-800 p-2 rounded-md hover:bg-blue-50 transition-colors"
                      title="Concluir Action Item"
                    >
                      <CheckCircle className="h-4 w-4" />
                      <span className="text-xs text-gray-400">Concluir</span>
                    </button>
                  )}

                  {/* View Feedback Button - only for done status */}
                  {actionItem.status === 'done' && extractFeedback(actionItem.description) && (
                    <button
                      onClick={() => handleViewFeedback(actionItem)}
                      className="flex flex-col items-center space-y-1 text-green-600 hover:text-green-800 p-2 rounded-md hover:bg-green-50 transition-colors"
                      title="Ver Parecer"
                    >
                      <FileText className="h-4 w-4" />
                      <span className="text-xs text-gray-400">Parecer</span>
                    </button>
                  )}

                        {/* Delete Button - only for todo and in_progress */}
                        {(actionItem.status === 'todo' || actionItem.status === 'in_progress') && canEditActionItem(actionItem) && (
                          <button
                            onClick={() => handleDeleteActionItem(actionItem)}
                            className="flex flex-col items-center space-y-1 text-red-600 hover:text-red-800 p-2 rounded-md hover:bg-red-50 transition-colors"
                            title="Excluir Action Item"
                          >
                            <Trash2 className="h-4 w-4" />
                            <span className="text-xs text-gray-400">Excluir</span>
                          </button>
                        )}

                  {/* View Retrospective Link */}
                  <Link
                    to={`/retrospectives/${actionItem.retrospective.id}`}
                    className="flex flex-col items-center space-y-1 text-gray-600 hover:text-gray-900 p-2 rounded-md hover:bg-gray-50 transition-colors"
                    title="Ver Retrospectiva"
                  >
                    <MessageSquare className="h-4 w-4" />
                    <span className="text-xs text-gray-400">Ver</span>
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <CheckSquare className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            Nenhum action item encontrado
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            Tente ajustar os filtros para encontrar action items.
          </p>
        </div>
      )}

      {/* Complete Action Item Modal */}
      {showCompleteModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-6 w-11/12 md:w-2/3 lg:w-1/2 bg-white rounded-lg shadow-xl">
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-medium text-gray-900">
                ✅ Concluir Action Item - {completingActionItem?.title}
              </h3>
              <button
                onClick={() => setShowCompleteModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="h-6 w-6" />
              </button>
            </div>

            {/* Form */}
            <form onSubmit={handleSaveCompletion} className="space-y-4">
              {/* Feedback */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  📝 Parecer sobre a conclusão *
                </label>
                <textarea
                  value={completionForm.feedback}
                  onChange={(e) => setCompletionForm({ ...completionForm, feedback: e.target.value })}
                  rows={4}
                  className="input"
                  placeholder="Descreva como foi a execução, quais foram os resultados obtidos, lições aprendidas, etc..."
                  required
                />
              </div>

              {/* Retrospective Info */}
              {completingActionItem && (
                <div className="bg-gray-50 p-4 rounded-md">
                  <h4 className="text-sm font-medium text-gray-700 mb-2">Retrospectiva de Origem</h4>
                  <div className="flex items-center space-x-2 text-sm text-gray-600">
                    <MessageSquare className="h-4 w-4" />
                    <span>{completingActionItem.retrospective.title}</span>
                    <span className="text-gray-400">•</span>
                    <span>{getTemplateIcon(completingActionItem.retrospective.template)} {completingActionItem.retrospective.template}</span>
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="flex items-center justify-end space-x-3 pt-6">
                <button
                  type="button"
                  onClick={() => setShowCompleteModal(false)}
                  className="btn btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={updateActionItemMutation.isLoading}
                  className="btn btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {updateActionItemMutation.isLoading ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      Concluindo...
                    </>
                  ) : (
                    <>
                      <CheckCircle className="h-4 w-4 mr-2" />
                      Concluir Action Item
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* View Feedback Modal */}
      {showFeedbackModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-6 w-11/12 md:w-2/3 lg:w-1/2 bg-white rounded-lg shadow-xl">
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-medium text-gray-900">
                📝 Parecer - {viewingActionItem?.title}
              </h3>
              <button
                onClick={() => setShowFeedbackModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="h-6 w-6" />
              </button>
            </div>

            {/* Content */}
            <div className="space-y-4">
              {/* Feedback */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Parecer sobre a conclusão
                </label>
                <div className="bg-gray-50 p-4 rounded-md border">
                  <p className="text-gray-800 whitespace-pre-wrap">
                    {extractFeedback(viewingActionItem?.description) || 'Nenhum parecer disponível.'}
                  </p>
                </div>
              </div>

              {/* Retrospective Info */}
              {viewingActionItem && (
                <div className="bg-gray-50 p-4 rounded-md">
                  <h4 className="text-sm font-medium text-gray-700 mb-2">Retrospectiva de Origem</h4>
                  <div className="flex items-center space-x-2 text-sm text-gray-600">
                    <MessageSquare className="h-4 w-4" />
                    <span>{viewingActionItem.retrospective.title}</span>
                    <span className="text-gray-400">•</span>
                    <span>{getTemplateIcon(viewingActionItem.retrospective.template)} {viewingActionItem.retrospective.template}</span>
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="flex items-center justify-end space-x-3 pt-6">
                <button
                  type="button"
                  onClick={() => setShowFeedbackModal(false)}
                  className="btn btn-secondary"
                >
                  Fechar
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setDeletingActionItem(null);
        }}
        onConfirm={confirmDeleteActionItem}
        title="Excluir Action Item"
        message={`Tem certeza que deseja excluir o Action Item "${deletingActionItem?.title}"? Esta ação não pode ser desfeita.`}
        confirmText="Excluir"
        cancelText="Cancelar"
      />
    </div>
  );
};

export default ActionItemsPage;
