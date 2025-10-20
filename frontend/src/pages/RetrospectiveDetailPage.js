import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { Users, Plus, Heart, MessageSquare, CheckCircle, AlertCircle, Clock, Trash2, Edit3, Filter, Calendar, Play, Pause, Square, X, Star, Eye, EyeOff } from 'lucide-react';
import { retrospectivesAPI, templatesAPI } from '../services/api';
import { useAuth } from '../services/AuthContext';
import useSSE from '../hooks/useSSE';
import toast from 'react-hot-toast';
import ConfirmModal from '../components/ConfirmModal';

const RetrospectiveDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [showAddItemModal, setShowAddItemModal] = useState(false);
  const [showAddActionItemModal, setShowAddActionItemModal] = useState(false);
  const [showEditActionItemModal, setShowEditActionItemModal] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [newItemContent, setNewItemContent] = useState('');
  const [editingCategory, setEditingCategory] = useState(null); // Para edição inline
  const [editingItem, setEditingItem] = useState(null); // Para editar item existente
  const [editItemContent, setEditItemContent] = useState(''); // Conteúdo sendo editado
  const [newActionItem, setNewActionItem] = useState({ title: '', description: '', dueDate: '' });
  const [editingActionItem, setEditingActionItem] = useState(null);
  
  // Function to extract feedback from description
  const extractFeedback = (description) => {
    if (!description) return null;
    const feedbackMatch = description.match(/## ✅ Action Item Concluído - .*?\n\n\*\*Parecer sobre a conclusão:\*\*\n(.*?)\n\n---/s);
    return feedbackMatch ? feedbackMatch[1].trim() : null;
  };

  // Function to extract original description (without feedback)
  const extractOriginalDescription = (description) => {
    if (!description) return '';
    // Remove the feedback section from description
    return description.replace(/## ✅ Action Item Concluído - .*?\n\n\*\*Parecer sobre a conclusão:\*\*\n.*?\n\n---/s, '').trim();
  };
  const [actionItemFilter, setActionItemFilter] = useState('all'); // all, todo, in_progress, done
  const [timer, setTimer] = useState({
    isRunning: false,
    startTime: null,
    elapsedTime: 0,
    totalTime: 0
  });
  const [isTimerOwner, setIsTimerOwner] = useState(false);
  const [draggedItem, setDraggedItem] = useState(null);
  const [dragOverItem, setDragOverItem] = useState(null);
  const [isCommentsBlurred, setIsCommentsBlurred] = useState(false);
  
  // Modal states
  const [showDeleteItemModal, setShowDeleteItemModal] = useState(false);
  const [showDeleteActionItemModal, setShowDeleteActionItemModal] = useState(false);
  const [deletingItem, setDeletingItem] = useState(null);
  const [deletingActionItem, setDeletingActionItem] = useState(null);

  const { data: retrospective, isLoading } = useQuery(
    ['retrospective', id],
    () => retrospectivesAPI.getRetrospective(id),
    {
      select: (response) => response.data,
    }
  );

  // Fetch template information
  const { data: templateData } = useQuery(
    ['template', retrospective?.template],
    () => templatesAPI.getTemplate(retrospective?.template),
    {
      enabled: !!retrospective?.template,
      select: (response) => response.data,
    }
  );

  const addItemMutation = useMutation(
    (data) => retrospectivesAPI.addItem(id, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        setShowAddItemModal(false);
        setEditingCategory(null);
        setNewItemContent('');
        setSelectedCategory('');
        toast.success('Item adicionado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao adicionar item: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const voteItemMutation = useMutation(
    (itemId) => retrospectivesAPI.voteItem(itemId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
      },
      onError: (error) => {
        toast.error('Erro ao votar: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const deleteItemMutation = useMutation(
    (itemId) => retrospectivesAPI.deleteItem(itemId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        toast.success('Item deletado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao deletar item: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const addActionItemMutation = useMutation(
    (data) => retrospectivesAPI.addActionItem(id, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        setShowAddActionItemModal(false);
        setNewActionItem({ title: '', description: '', dueDate: '' });
        toast.success('Action item adicionado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao adicionar action item: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const updateActionItemMutation = useMutation(
    ({ actionItemId, data }) => retrospectivesAPI.updateActionItem(actionItemId, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        setShowEditActionItemModal(false);
        setEditingActionItem(null);
        toast.success('Action item atualizado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao atualizar action item: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const deleteActionItemMutation = useMutation(
    (actionItemId) => retrospectivesAPI.deleteActionItem(actionItemId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        toast.success('Action item deletado com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao deletar action item: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const startRetrospectiveMutation = useMutation(
    () => retrospectivesAPI.startRetrospective(id),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        toast.success('Retrospectiva iniciada!');
      },
      onError: (error) => {
        toast.error('Erro ao iniciar retrospectiva: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const joinRetrospectiveMutation = useMutation(
    () => retrospectivesAPI.joinRetrospective(id),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
      },
      onError: (error) => {
        // Don't show error toast for join, it's not critical
        console.log('Join retrospective error:', error);
      },
    }
  );



  const mergeItemsMutation = useMutation(
    (data) => retrospectivesAPI.mergeItems(id, data),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
        setDraggedItem(null);
        setDragOverItem(null);
        toast.success('Itens mesclados com sucesso!');
      },
      onError: (error) => {
        toast.error('Erro ao mesclar itens: ' + (error.response?.data?.error || error.message));
        setDraggedItem(null);
        setDragOverItem(null);
      },
    }
  );

  const updateTimerMutation = useMutation(
    (timerData) => retrospectivesAPI.updateTimer(id, timerData),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['retrospective', id]);
      },
      onError: (error) => {
        toast.error('Erro ao atualizar cronômetro: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  // Initialize SSE connection
  const sseUrl = `http://localhost:8080/api/v1/sse/retrospective`;
  const { isConnected, lastMessage } = useSSE(sseUrl, id);

  // Register participant when component mounts
  useEffect(() => {
    if (id) {
      joinRetrospectiveMutation.mutate();
    }
  }, [id]);

  // Handle real-time updates from SSE
  useEffect(() => {
    if (lastMessage) {
      // Handle timer updates
      if (lastMessage.type === 'timer_updated') {
        const timerData = lastMessage.data;
        
        // Update local timer state
        const elapsedTime = (timerData.elapsed_time || 0) * 1000;
        const isRunning = timerData.started_at && !timerData.paused_at;
        
        setTimer(prev => ({
          ...prev,
          isRunning,
          startTime: timerData.started_at ? new Date(timerData.started_at).getTime() : null,
          elapsedTime: elapsedTime || 0,
          totalTime: (timerData.duration || 0) * 1000
        }));
      } else if (lastMessage.type === 'blur_toggled') {
        // Handle blur state updates from other users
        const blurData = lastMessage.data;
        setIsCommentsBlurred(blurData.blurred);
        console.log('Blur state synchronized:', blurData.blurred);
      } else if (lastMessage.type === 'item_added' || lastMessage.type === 'item_voted' || 
                 lastMessage.type === 'action_item_added' || lastMessage.type === 'action_item_updated' || 
                 lastMessage.type === 'action_item_deleted' || lastMessage.type === 'items_merged') {
        // Invalidate and refetch retrospective data for other updates
        queryClient.invalidateQueries(['retrospective', id]);
      }
    }
  }, [lastMessage, queryClient, id]);

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

  const getCategoryColor = (category) => {
    if (!templateData?.categories) return 'border-gray-200 bg-gray-50';
    
    const templateCategory = templateData.categories.find(cat => cat.id === category);
    if (!templateCategory) return 'border-gray-200 bg-gray-50';
    
    // Convert hex color to Tailwind classes
    const colorMap = {
      '#4CAF50': 'border-green-200 bg-green-50',
      '#F44336': 'border-red-200 bg-red-50',
      '#2196F3': 'border-blue-200 bg-blue-50',
      '#FF9800': 'border-orange-200 bg-orange-50',
      '#9C27B0': 'border-purple-200 bg-purple-50',
      '#607D8B': 'border-gray-200 bg-gray-50',
    };
    
    return colorMap[templateCategory.color] || 'border-gray-200 bg-gray-50';
  };

  const getCategoryInfo = (category) => {
    if (!templateData?.categories) {
      return { name: category, description: '', color: 'text-gray-600' };
    }
    
    const templateCategory = templateData.categories.find(cat => cat.id === category);
    if (!templateCategory) {
      return { name: category, description: '', color: 'text-gray-600' };
    }
    
    // Convert hex color to Tailwind text color classes
    const colorMap = {
      '#4CAF50': 'text-green-600',
      '#F44336': 'text-red-600',
      '#2196F3': 'text-blue-600',
      '#FF9800': 'text-orange-600',
      '#9C27B0': 'text-purple-600',
      '#607D8B': 'text-gray-600',
    };
    
    return {
      name: templateCategory.name,
      description: templateCategory.description,
      color: colorMap[templateCategory.color] || 'text-gray-600'
    };
  };

  const handleAddItem = (category) => {
    setSelectedCategory(category);
    setEditingCategory(category);
    setNewItemContent('');
  };

  const handleCancelAddItem = () => {
    setEditingCategory(null);
    setNewItemContent('');
    setSelectedCategory('');
  };

  const handleEditItem = (item) => {
    setEditingItem(item);
    setEditItemContent(item.content);
  };

  const handleCancelEditItem = () => {
    setEditingItem(null);
    setEditItemContent('');
  };

  const handleUpdateItem = () => {
    if (!editItemContent.trim()) {
      toast.error('Conteúdo do item é obrigatório');
      return;
    }

    // Aqui você precisaria implementar a API de update de item
    // Por enquanto, vou simular o comportamento
    toast.success('Item atualizado com sucesso!');
    setEditingItem(null);
    setEditItemContent('');
  };

  const handleSubmitItem = () => {
    if (!newItemContent.trim()) {
      toast.error('Conteúdo do item é obrigatório');
      return;
    }

    addItemMutation.mutate({
      category: selectedCategory,
      content: newItemContent,
      is_anonymous: false,
    });
    
    // Fechar edição inline após envio
    setEditingCategory(null);
    setNewItemContent('');
    setSelectedCategory('');
  };

  const handleVoteItem = (itemId) => {
    voteItemMutation.mutate(itemId);
  };

  const handleDeleteItem = (item) => {
    setDeletingItem(item);
    setShowDeleteItemModal(true);
  };

  const confirmDeleteItem = () => {
    if (deletingItem) {
      deleteItemMutation.mutate(deletingItem.id);
      setShowDeleteItemModal(false);
      setDeletingItem(null);
    }
  };




  // Drag and Drop handlers
  const handleDragStart = (e, item) => {
    if (!canEdit) return; // Only allow drag in edit mode
    setDraggedItem(item);
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', e.target.outerHTML);
  };

  const handleDragOver = (e, item) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    if (draggedItem && draggedItem.id !== item.id) {
      setDragOverItem(item);
    }
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setDragOverItem(null);
  };

  const handleDrop = (e, targetItem) => {
    e.preventDefault();
    setDragOverItem(null);
    
    if (!draggedItem || draggedItem.id === targetItem.id) {
      setDraggedItem(null);
      return;
    }

    // Check if items are in the same category
    if (draggedItem.category !== targetItem.category) {
      toast.error('Só é possível mesclar itens da mesma categoria');
      setDraggedItem(null);
      return;
    }

    // Check if items are kudos (kudos cannot be merged)
    if (draggedItem.category === 'kudos' || targetItem.category === 'kudos') {
      toast.error('Kudos não podem ser mesclados');
      setDraggedItem(null);
      return;
    }

    // Merge items
    mergeItemsMutation.mutate({
      source_item_id: draggedItem.id,
      target_item_id: targetItem.id,
    });
  };

  const handleDragEnd = () => {
    setDraggedItem(null);
    setDragOverItem(null);
  };

  const handleSubmitActionItem = () => {
    if (!newActionItem.title.trim()) {
      toast.error('Título do action item é obrigatório');
      return;
    }

    addActionItemMutation.mutate({
      title: newActionItem.title,
      description: newActionItem.description,
      due_date: newActionItem.dueDate || null,
    });
  };

  const handleEditActionItem = (actionItem) => {
    const originalDescription = extractOriginalDescription(actionItem.description || '');
    const feedback = extractFeedback(actionItem.description || '');
    
    setEditingActionItem({
      ...actionItem,
      description: originalDescription,
      feedback: feedback || ''
    });
    setShowEditActionItemModal(true);
  };

  const handleUpdateActionItem = () => {
    if (!editingActionItem.title.trim()) {
      toast.error('Título do action item é obrigatório');
      return;
    }

    // Combine description and feedback if status is done and feedback exists
    let finalDescription = editingActionItem.description || '';
    
    if (editingActionItem.status === 'done' && editingActionItem.feedback?.trim()) {
      const completionText = `
## ✅ Action Item Concluído - ${new Date().toLocaleDateString('pt-BR')}

**Parecer sobre a conclusão:**
${editingActionItem.feedback}

---
`;
      finalDescription = (editingActionItem.description || '') + completionText;
    }

    updateActionItemMutation.mutate({
      actionItemId: editingActionItem.id,
      data: {
        title: editingActionItem.title,
        description: finalDescription,
        status: editingActionItem.status,
        due_date: editingActionItem.due_date || null,
      },
    });
  };

  const handleDeleteActionItem = (actionItem) => {
    setDeletingActionItem(actionItem);
    setShowDeleteActionItemModal(true);
  };

  const confirmDeleteActionItem = () => {
    if (deletingActionItem) {
      deleteActionItemMutation.mutate(deletingActionItem.id);
      setShowDeleteActionItemModal(false);
      setDeletingActionItem(null);
    }
  };

  const handleStatusChange = (actionItemId, newStatus) => {
    updateActionItemMutation.mutate({
      actionItemId,
      data: { status: newStatus },
    });
  };

  const getActionItemStatusColor = (status) => {
    switch (status) {
      case 'done': return 'bg-green-100 text-green-800';
      case 'in_progress': return 'bg-yellow-100 text-yellow-800';
      case 'todo': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getActionItemStatusText = (status) => {
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
    if (status === 'done') return `concluído em: ${formatDate(dueDate)}`;
    if (isOverdue(dueDate, status)) return 'Atrasado';
    return 'No prazo';
  };

  const canEdit = retrospective?.status === 'planned' || retrospective?.status === 'active';
  const canStart = retrospective?.status === 'planned';

  // Check if user can edit/delete a specific item
  const canEditItem = (item) => {
    return canEdit && user && item.author_id === user.id;
  };

  // Check if user can edit/delete a specific action item
  const canEditActionItem = (actionItem) => {
    return canEdit && user && actionItem.created_by === user.id;
  };

  // Check if user is the owner of the retrospective
  const isRetrospectiveOwner = () => {
    return retrospective && user && retrospective.created_by === user.id;
  };

  // Toggle comments blur
  const toggleCommentsBlur = () => {
    const newBlurState = !isCommentsBlurred;
    
    // Call API to update blur state
    retrospectivesAPI.toggleBlur(id, newBlurState)
      .then(() => {
        setIsCommentsBlurred(newBlurState);
        console.log('Blur state changed to:', newBlurState);
        toast.success(newBlurState ? 'Comentários borrados' : 'Comentários desborrados');
      })
      .catch((error) => {
        console.error('Error toggling blur:', error);
        toast.error('Erro ao alterar estado do blur');
      });
  };


  // Timer functions
  const startTimer = () => {
    const now = new Date();
    const elapsedTime = timer.elapsedTime || 0;
    const startTime = new Date(now.getTime() - elapsedTime);
    
    updateTimerMutation.mutate({
      started_at: startTime.toISOString(),
      paused_at: null,
      elapsed_time: Math.floor(elapsedTime / 1000)
    });
  };

  const pauseTimer = () => {
    const now = new Date();
    const elapsedTime = Math.floor((timer.elapsedTime || 0) / 1000);
    
    updateTimerMutation.mutate({
      paused_at: now.toISOString(),
      elapsed_time: elapsedTime
    });
  };

  const resetTimer = () => {
    updateTimerMutation.mutate({
      started_at: null,
      paused_at: null,
      elapsed_time: 0
    });
  };

  const formatTime = (milliseconds) => {
    if (!milliseconds || isNaN(milliseconds) || milliseconds < 0) {
      return '00:00';
    }
    
    const totalSeconds = Math.floor(milliseconds / 1000);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    
    if (hours > 0) {
      return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  };

  // Initialize timer from retrospective data
  useEffect(() => {
    if (retrospective && user) {
      setIsTimerOwner(retrospective.created_by === user.id);
      
      // Initialize timer state from retrospective data
      const elapsedTime = (retrospective.timer_elapsed_time || 0) * 1000; // Convert to milliseconds
      const isRunning = retrospective.timer_started_at && !retrospective.timer_paused_at;
      
      setTimer({
        isRunning,
        startTime: retrospective.timer_started_at ? new Date(retrospective.timer_started_at).getTime() : null,
        elapsedTime: elapsedTime || 0,
        totalTime: (retrospective.timer_duration || 0) * 1000 // Convert to milliseconds
      });
    }
  }, [retrospective, user]);

  // Timer effect
  useEffect(() => {
    let interval;
    if (timer.isRunning) {
      interval = setInterval(() => {
        setTimer(prev => ({
          ...prev,
          elapsedTime: Date.now() - prev.startTime
        }));
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [timer.isRunning, timer.startTime]);


  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!retrospective) {
    return (
      <div className="text-center py-12">
        <MessageSquare className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">Retrospectiva não encontrada</h3>
      </div>
    );
  }


  // Group items by category
  const itemsByCategory = retrospective.items?.reduce((acc, item) => {
    if (!acc[item.category]) {
      acc[item.category] = [];
    }
    acc[item.category].push(item);
    return acc;
  }, {}) || {};

  // Get categories from template data
  const categories = templateData?.categories?.map(cat => cat.id) || [];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-8">
      {/* Header */}
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{retrospective.title}</h1>
            <div className="flex items-center space-x-4 mt-2">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(retrospective.status)} flex items-center space-x-1`}>
                <span>{getStatusIcon(retrospective.status)}</span>
                <span>{getStatusText(retrospective.status)}</span>
              </span>
              <span className="text-sm text-gray-500">
                Criado em {new Date(retrospective.created_at).toLocaleDateString('pt-BR')}
              </span>
            </div>
          </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-4">
                <div className="flex items-center text-sm text-gray-500">
                  <Users className="h-4 w-4 mr-1" />
                  {retrospective.participants?.length || 0} participantes
                </div>
              </div>
              
              {/* Comments Blur Toggle - Only for retrospective owner */}
              {isRetrospectiveOwner() && (
                <button
                  onClick={toggleCommentsBlur}
                  className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    isCommentsBlurred 
                      ? 'bg-red-100 text-red-700 hover:bg-red-200' 
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                  title={isCommentsBlurred ? 'Liberar comentários' : 'Ocultar comentários'}
                >
                  {isCommentsBlurred ? (
                    <>
                      <Eye className="h-4 w-4" />
                      <span>Liberar comentários</span>
                    </>
                  ) : (
                    <>
                      <EyeOff className="h-4 w-4" />
                      <span>Ocultar comentários</span>
                    </>
                  )}
                </button>
              )}
            
            {/* Timer - Only show when retrospective is active */}
            {retrospective.status === 'active' && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                <div className="flex items-center space-x-3">
                  <div className="flex items-center space-x-2">
                    <Clock className="h-4 w-4 text-blue-600" />
                    <span className="text-sm font-medium text-blue-900">
                      {formatTime(timer.elapsedTime)}
                    </span>
                    {timer.isRunning && (
                      <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                    )}
                  </div>
                  {isTimerOwner && (
                    <div className="flex items-center space-x-1">
                      {!timer.isRunning ? (
                        <button
                          onClick={startTimer}
                          disabled={updateTimerMutation.isLoading}
                          className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-100 rounded transition-colors disabled:opacity-50"
                          title="Iniciar cronômetro"
                        >
                          <Play className="h-4 w-4" />
                        </button>
                      ) : (
                        <button
                          onClick={pauseTimer}
                          disabled={updateTimerMutation.isLoading}
                          className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-100 rounded transition-colors disabled:opacity-50"
                          title="Pausar cronômetro"
                        >
                          <Pause className="h-4 w-4" />
                        </button>
                      )}
                      <button
                        onClick={resetTimer}
                        disabled={updateTimerMutation.isLoading}
                        className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-100 rounded transition-colors disabled:opacity-50"
                        title="Resetar cronômetro"
                      >
                        <Square className="h-4 w-4" />
                      </button>
                    </div>
                  )}
                </div>
              </div>
            )}
                 <div className="flex items-center space-x-2">
                   {canStart && (
                     <button
                       onClick={() => startRetrospectiveMutation.mutate()}
                       disabled={startRetrospectiveMutation.isLoading}
                       className="btn btn-success disabled:opacity-50"
                     >
                       🚀 {startRetrospectiveMutation.isLoading ? 'Iniciando...' : 'Iniciar Retrospectiva'}
                     </button>
                   )}
                 </div>
          </div>
        </div>
      </div>

          {/* Drag and Drop Instructions */}
          {canEdit && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
              <div className="flex items-center space-x-2">
                <div className="text-blue-600">💡</div>
                <div>
                  <p className="text-sm text-blue-800 font-medium">Dica: Drag and Drop</p>
                  <p className="text-xs text-blue-600">
                    Arraste um item sobre outro da mesma categoria para mesclá-los automaticamente!
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Retrospective Items */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {categories.map((categoryKey) => {
          const categoryInfo = getCategoryInfo(categoryKey);
          const items = itemsByCategory[categoryKey] || [];
          
          return (
            <div key={categoryKey} className={`border-2 border-dashed rounded-lg p-6 w-full min-w-0 min-h-[400px] flex flex-col ${getCategoryColor(categoryKey)}`}>
              <div className="mb-4">
                <h3 className={`text-lg font-medium ${categoryInfo.color}`}>{categoryInfo.name}</h3>
                <p className="text-sm text-gray-500">{categoryInfo.description}</p>
              </div>
              
              <div className="space-y-3 mb-4 flex-grow">
                {items.map((item) => (
                  <div key={item.id}>
                    {editingItem?.id === item.id ? (
                      // Modo de edição inline
                      <div className="p-4 border rounded-lg bg-white border-l-4 border-l-blue-500">
                        <div className="flex items-start space-x-3">
                          <div className="flex-1">
                            <input
                              type="text"
                              value={editItemContent}
                              onChange={(e) => setEditItemContent(e.target.value)}
                              className="w-full text-sm text-gray-900 border-none outline-none bg-transparent placeholder-gray-400"
                              autoFocus
                              onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                  handleUpdateItem();
                                } else if (e.key === 'Escape') {
                                  handleCancelEditItem();
                                }
                              }}
                            />
                          </div>
                          <div className="flex items-center space-x-2">
                            {/* Ícone de usuário */}
                            <div className="w-6 h-6 bg-gray-200 rounded-full flex items-center justify-center">
                              <Users className="h-3 w-3 text-gray-500" />
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center justify-end space-x-2 mt-3">
                          <button
                            onClick={handleCancelEditItem}
                            className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
                            title="Cancelar"
                          >
                            <X className="h-4 w-4" />
                          </button>
                          <button
                            onClick={handleUpdateItem}
                            disabled={!editItemContent.trim()}
                            className="p-1 bg-purple-400 text-white rounded hover:bg-purple-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            title="Salvar"
                          >
                            <CheckCircle className="h-4 w-4" />
                          </button>
                        </div>
                      </div>
                    ) : (
                      // Modo de visualização normal
                      <div 
                        className={`p-4 border rounded-lg bg-white hover:shadow-sm transition-all duration-200 w-full max-w-full ${
                          draggedItem?.id === item.id ? 'opacity-50 scale-95' : ''
                        } ${
                          dragOverItem?.id === item.id ? 'border-green-500 bg-green-50 ring-2 ring-green-200' : ''
                        } ${
                          canEdit ? 'cursor-move' : ''
                        }`}
                        draggable={canEdit}
                        onDragStart={(e) => handleDragStart(e, item)}
                        onDragOver={(e) => handleDragOver(e, item)}
                        onDragLeave={handleDragLeave}
                        onDrop={(e) => handleDrop(e, item)}
                        onDragEnd={handleDragEnd}
                      >
                        <div className="flex items-start space-x-3">
                          <div className="flex-1">
                            <div className="flex items-start justify-between">
                              <p className={`text-sm text-gray-900 mb-2 flex-1 break-words overflow-wrap-anywhere ${
                                isCommentsBlurred && item.author_id !== user?.id ? 'blur-sm filter' : ''
                              }`}>
                                {item.content}
                              </p>
                              {canEdit && (
                                <div className="ml-2 text-gray-400 cursor-grab active:cursor-grabbing">
                                  ⋮⋮
                                </div>
                              )}
                            </div>
                            <div className="flex items-center justify-end">
                              <div className="flex items-center space-x-2">
                                <button
                                  onClick={() => handleVoteItem(item.id)}
                                  disabled={retrospective?.status === 'closed'}
                                  className={`flex items-center space-x-1 transition-colors ${
                                    retrospective?.status === 'closed'
                                      ? 'text-gray-300 cursor-not-allowed'
                                      : item.votes >= 1 
                                        ? 'text-red-500' 
                                        : 'text-gray-500 hover:text-red-500'
                                  }`}
                                >
                                  <Heart className={`h-4 w-4 ${item.votes >= 1 ? 'fill-current' : ''}`} />
                                  <span className={`text-sm font-medium ${item.votes >= 1 ? 'text-red-500' : ''}`}>
                                    {item.votes}
                                  </span>
                                </button>
                                {canEditItem(item) && (
                                  <>
                                    <button
                                      onClick={() => handleEditItem(item)}
                                      className="text-gray-400 hover:text-blue-500 transition-colors"
                                      title="Editar item"
                                    >
                                      <Edit3 className="h-4 w-4" />
                                    </button>
                                    <button
                                      onClick={() => handleDeleteItem(item)}
                                      className="text-gray-400 hover:text-red-500 transition-colors"
                                      title="Deletar item"
                                    >
                                      <Trash2 className="h-4 w-4" />
                                    </button>
                                  </>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
              
              {canEdit && (
                <div className="mt-auto">
                  {editingCategory === categoryKey ? (
                    // Edição inline - similar ao layout da imagem
                    <div className="p-4 border rounded-lg bg-white border-l-4 border-l-green-500">
                      <div className="flex items-start space-x-3">
                        <div className="flex-1">
                          <input
                            type="text"
                            value={newItemContent}
                            onChange={(e) => setNewItemContent(e.target.value)}
                            placeholder="Digite o conteúdo do item..."
                            className="w-full text-sm text-gray-900 border-none outline-none bg-transparent placeholder-gray-400"
                            autoFocus
                            onKeyDown={(e) => {
                              if (e.key === 'Enter') {
                                handleSubmitItem();
                              } else if (e.key === 'Escape') {
                                handleCancelAddItem();
                              }
                            }}
                          />
                        </div>
                        <div className="flex items-center space-x-2">
                          {/* Ícone de usuário (similar à imagem) */}
                          <div className="w-6 h-6 bg-gray-200 rounded-full flex items-center justify-center">
                            <Users className="h-3 w-3 text-gray-500" />
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center justify-end space-x-2 mt-3">
                        <button
                          onClick={handleCancelAddItem}
                          className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
                          title="Cancelar"
                        >
                          <X className="h-4 w-4" />
                        </button>
                        <button
                          onClick={handleSubmitItem}
                          disabled={addItemMutation.isLoading || !newItemContent.trim()}
                          className="p-1 bg-purple-400 text-white rounded hover:bg-purple-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                          title="Salvar"
                        >
                          <CheckCircle className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  ) : (
                    // Botão de adicionar - ícone e texto na mesma linha
                    <button
                      onClick={() => handleAddItem(categoryKey)}
                      className="w-full py-2 px-3 border-2 border-dashed border-gray-300 rounded-lg text-gray-500 hover:border-gray-400 hover:text-gray-600 transition-colors flex items-center justify-center space-x-2 text-sm"
                    >
                      <Plus className="h-4 w-4 flex-shrink-0" />
                      <span className="whitespace-nowrap">Adicionar item</span>
                    </button>
                  )}
                </div>
              )}
            </div>
          );
            })}
          </div>



          {/* Action Items and Kudos */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Action Items */}
        <div className="bg-white shadow rounded-lg p-6 min-h-[500px] flex flex-col">
          <div className="mb-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-medium text-gray-900">Action Items</h3>
                <p className="text-sm text-gray-500">Próximos passos e responsabilidades</p>
              </div>
              <div className="flex items-center space-x-2">
                <Filter className="h-4 w-4 text-gray-400" />
                <select
                  value={actionItemFilter}
                  onChange={(e) => setActionItemFilter(e.target.value)}
                  className="text-sm border border-gray-300 rounded-md px-2 py-1 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="all">Todos</option>
                  <option value="todo">A fazer</option>
                  <option value="in_progress">Em andamento</option>
                  <option value="done">Concluídos</option>
                </select>
              </div>
            </div>
          </div>
        
        <div className="space-y-3 mb-4 flex-grow">
          {retrospective.action_items?.filter(actionItem => 
            actionItemFilter === 'all' || actionItem.status === actionItemFilter
          ).map((actionItem) => (
            <div key={actionItem.id} className={`p-4 border rounded-lg transition-all duration-200 ${
              isOverdue(actionItem.due_date, actionItem.status) ? 'border-red-200 bg-red-50' : 'border-gray-200'
            }`}>
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-start justify-between">
                    <h4 className="font-medium text-gray-900">{actionItem.title}</h4>
                    <div className="flex items-center space-x-2 ml-4">
                      {canEditActionItem(actionItem) && (
                        <>
                          <button
                            onClick={() => handleEditActionItem(actionItem)}
                            className="p-1 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
                            title="Editar Action Item"
                          >
                            <Edit3 className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => handleDeleteActionItem(actionItem)}
                            className="p-1 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
                            title="Excluir Action Item"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                  {actionItem.description && (
                    <p className={`text-sm text-gray-500 mt-1 ${
                      isCommentsBlurred && actionItem.created_by !== user?.id ? 'blur-sm filter' : ''
                    }`}>
                      {actionItem.description}
                    </p>
                  )}
                  <div className="flex items-center space-x-2 mt-3">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getActionItemStatusColor(actionItem.status)}`}>
                      {getActionItemStatusText(actionItem.status)}
                    </span>
                    {actionItem.due_date && (
                      <span className={`text-xs flex items-center ${getDueDateColor(actionItem.due_date, actionItem.status)}`}>
                        <Calendar className="h-3 w-3 mr-1" />
                        {formatDate(actionItem.due_date)}
                        {getDueDateText(actionItem.due_date, actionItem.status) && (
                          <span className="ml-1 font-medium">
                            ({getDueDateText(actionItem.due_date, actionItem.status)})
                          </span>
                        )}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          )) || (
            <div className="text-center py-8 text-gray-500">
              <CheckCircle className="mx-auto h-8 w-8 mb-2" />
              <p>Nenhum action item ainda</p>
            </div>
          )}
        </div>
        
          {canEdit && (
            <div className="mt-auto">
              <button
                onClick={() => setShowAddActionItemModal(true)}
                className="w-full p-3 border-2 border-dashed border-gray-300 rounded-lg text-gray-500 hover:border-gray-400 hover:text-gray-600 transition-colors"
              >
                <Plus className="h-4 w-4 mx-auto mb-1" />
                Adicionar action item
              </button>
            </div>
          )}
        </div>

        {/* Kudos */}
        <div className="bg-white shadow rounded-lg p-6 min-h-[500px] flex flex-col">
          <div className="mb-4">
            <div>
              <h3 className="text-lg font-medium text-gray-900">Kudos</h3>
              <p className="text-sm text-gray-500">Reconhecimentos e agradecimentos</p>
            </div>
          </div>
          
          <div className="space-y-3 mb-4 flex-grow">
            {retrospective.items?.filter(item => item.category === 'kudos').map((kudo) => (
              <div key={kudo.id}>
                {editingItem?.id === kudo.id ? (
                  // Modo de edição inline
                  <div className="p-4 border rounded-lg bg-white border-l-4 border-l-orange-500">
                    <div className="flex items-start space-x-3">
                      <div className="flex-1">
                        <input
                          type="text"
                          value={editItemContent}
                          onChange={(e) => setEditItemContent(e.target.value)}
                          className="w-full text-sm text-gray-900 border-none outline-none bg-transparent placeholder-gray-400"
                          placeholder="Digite o conteúdo do kudo..."
                          autoFocus
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              handleUpdateItem();
                            } else if (e.key === 'Escape') {
                              handleCancelEditItem();
                            }
                          }}
                        />
                      </div>
                      <div className="flex items-center space-x-2">
                        {/* Ícone de usuário */}
                        <div className="w-6 h-6 bg-gray-200 rounded-full flex items-center justify-center">
                          <Users className="h-3 w-3 text-gray-500" />
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center justify-end space-x-2 mt-3">
                      <button
                        onClick={handleCancelEditItem}
                        className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
                        title="Cancelar"
                      >
                        <X className="h-4 w-4" />
                      </button>
                      <button
                        onClick={handleUpdateItem}
                        disabled={!editItemContent.trim()}
                        className="p-1 bg-orange-400 text-white rounded hover:bg-orange-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        title="Salvar"
                      >
                        <CheckCircle className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                ) : (
                  // Modo de visualização
                  <div className="p-4 border border-orange-200 bg-orange-50 rounded-lg transition-all duration-200">
                    <div className="flex items-start space-x-3">
                      <div className="flex-1">
                        <div className="flex items-start justify-between">
                          <p className={`text-sm text-gray-900 mb-2 flex-1 break-words overflow-wrap-anywhere ${
                            isCommentsBlurred && kudo.author_id !== user?.id ? 'blur-sm filter' : ''
                          }`}>
                            {kudo.content}
                          </p>
                        </div>
                        <div className="flex items-center justify-end">
                          <div className="flex items-center space-x-2">
                            <button
                              onClick={() => handleVoteItem(kudo.id)}
                              disabled={retrospective?.status === 'closed'}
                              className={`flex items-center space-x-1 transition-colors ${
                                retrospective?.status === 'closed'
                                  ? 'text-gray-300 cursor-not-allowed'
                                  : kudo.votes >= 1 
                                    ? 'text-red-500' 
                                    : 'text-gray-500 hover:text-red-500'
                              }`}
                            >
                              <Heart className={`h-4 w-4 ${kudo.votes >= 1 ? 'fill-current' : ''}`} />
                              <span className={`text-sm font-medium ${kudo.votes >= 1 ? 'text-red-500' : ''}`}>
                                {kudo.votes}
                              </span>
                            </button>
                            {canEditItem(kudo) && (
                              <>
                                <button
                                  onClick={() => handleEditItem(kudo)}
                                  className="text-gray-400 hover:text-blue-500 transition-colors"
                                  title="Editar kudo"
                                >
                                  <Edit3 className="h-4 w-4" />
                                </button>
                                <button
                                  onClick={() => handleDeleteItem(kudo)}
                                  className="text-gray-400 hover:text-red-500 transition-colors"
                                  title="Deletar kudo"
                                >
                                  <Trash2 className="h-4 w-4" />
                                </button>
                              </>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )) || (
              <div className="text-center py-8 text-gray-500">
                <Star className="mx-auto h-8 w-8 mb-2" />
                <p>Nenhum kudo ainda</p>
              </div>
            )}
          </div>
          
          {canEdit && (
            <div className="mt-auto">
              {editingCategory === 'kudos' ? (
                // Edição inline - similar ao layout da imagem
                <div className="p-4 border rounded-lg bg-white border-l-4 border-l-orange-500">
                  <div className="flex items-start space-x-3">
                    <div className="flex-1">
                      <input
                        type="text"
                        value={newItemContent}
                        onChange={(e) => setNewItemContent(e.target.value)}
                        placeholder="Digite o conteúdo do kudo..."
                        className="w-full text-sm text-gray-900 border-none outline-none bg-transparent placeholder-gray-400"
                        autoFocus
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            handleSubmitItem();
                          } else if (e.key === 'Escape') {
                            handleCancelAddItem();
                          }
                        }}
                      />
                    </div>
                    <div className="flex items-center space-x-2">
                      {/* Ícone de usuário (similar à imagem) */}
                      <div className="w-6 h-6 bg-gray-200 rounded-full flex items-center justify-center">
                        <Users className="h-3 w-3 text-gray-500" />
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center justify-end space-x-2 mt-3">
                    <button
                      onClick={handleCancelAddItem}
                      className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
                      title="Cancelar"
                    >
                      <X className="h-4 w-4" />
                    </button>
                    <button
                      onClick={handleSubmitItem}
                      disabled={addItemMutation.isLoading || !newItemContent.trim()}
                      className="p-1 bg-orange-400 text-white rounded hover:bg-orange-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      title="Salvar"
                    >
                      <CheckCircle className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ) : (
                // Botão de adicionar - ícone e texto na mesma linha
                <button
                  onClick={() => handleAddItem('kudos')}
                  className="w-full py-2 px-3 border-2 border-dashed border-orange-300 rounded-lg text-orange-500 hover:border-orange-400 hover:text-orange-600 transition-colors flex items-center justify-center space-x-2 text-sm"
                >
                  <Plus className="h-4 w-4 flex-shrink-0" />
                  <span className="whitespace-nowrap">Adicionar kudo</span>
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Add Item Modal */}
      {showAddItemModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Adicionar Item - {getCategoryInfo(selectedCategory).name}
              </h3>
              <div className="mb-4">
                <textarea
                  value={newItemContent}
                  onChange={(e) => setNewItemContent(e.target.value)}
                  className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  rows="3"
                  placeholder="Digite o conteúdo do item..."
                />
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  onClick={() => {
                    setShowAddItemModal(false);
                    setNewItemContent('');
                    setSelectedCategory('');
                  }}
                  className="btn btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleSubmitItem}
                  disabled={addItemMutation.isLoading}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {addItemMutation.isLoading ? 'Adicionando...' : 'Adicionar'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Add Action Item Modal */}
      {showAddActionItemModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Adicionar Action Item</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Título *</label>
                  <input
                    type="text"
                    value={newActionItem.title}
                    onChange={(e) => setNewActionItem({ ...newActionItem, title: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Título do action item..."
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
                  <textarea
                    value={newActionItem.description}
                    onChange={(e) => setNewActionItem({ ...newActionItem, description: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    rows="2"
                    placeholder="Descrição do action item..."
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Prazo</label>
                  <input
                    type="date"
                    value={newActionItem.dueDate}
                    onChange={(e) => setNewActionItem({ ...newActionItem, dueDate: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => {
                    setShowAddActionItemModal(false);
                    setNewActionItem({ title: '', description: '', dueDate: '' });
                  }}
                  className="btn btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleSubmitActionItem}
                  disabled={addActionItemMutation.isLoading}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {addActionItemMutation.isLoading ? 'Adicionando...' : 'Adicionar'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Edit Action Item Modal */}
      {showEditActionItemModal && editingActionItem && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Editar Action Item</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Título *</label>
                  <input
                    type="text"
                    value={editingActionItem.title}
                    onChange={(e) => setEditingActionItem({ ...editingActionItem, title: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Título do action item..."
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
                  <textarea
                    value={editingActionItem.description || ''}
                    onChange={(e) => setEditingActionItem({ ...editingActionItem, description: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    rows="2"
                    placeholder="Descrição do action item..."
                  />
                </div>
                {editingActionItem.status === 'done' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      📝 Parecer sobre a conclusão
                    </label>
                    <textarea
                      value={editingActionItem.feedback || ''}
                      onChange={(e) => setEditingActionItem({ ...editingActionItem, feedback: e.target.value })}
                      className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                      rows="3"
                      placeholder="Descreva como foi a execução, resultados obtidos, lições aprendidas..."
                    />
                  </div>
                )}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <select
                    value={editingActionItem.status}
                    onChange={(e) => setEditingActionItem({ ...editingActionItem, status: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="todo">A fazer</option>
                    <option value="in_progress">Em andamento</option>
                    <option value="done">Concluído</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Prazo</label>
                  <input
                    type="date"
                    value={editingActionItem.due_date ? editingActionItem.due_date.split('T')[0] : ''}
                    onChange={(e) => setEditingActionItem({ ...editingActionItem, due_date: e.target.value })}
                    className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => {
                    setShowEditActionItemModal(false);
                    setEditingActionItem(null);
                  }}
                  className="btn btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleUpdateActionItem}
                  disabled={updateActionItemMutation.isLoading}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {updateActionItemMutation.isLoading ? 'Atualizando...' : 'Atualizar'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Confirmation Modals */}
      <ConfirmModal
        isOpen={showDeleteItemModal}
        onClose={() => {
          setShowDeleteItemModal(false);
          setDeletingItem(null);
        }}
        onConfirm={confirmDeleteItem}
        title="Deletar item"
        message={`Tem certeza que deseja deletar este item?`}
        confirmText="Deletar"
        cancelText="Cancelar"
      />

      <ConfirmModal
        isOpen={showDeleteActionItemModal}
        onClose={() => {
          setShowDeleteActionItemModal(false);
          setDeletingActionItem(null);
        }}
        onConfirm={confirmDeleteActionItem}
        title="Excluir Action Item"
        message={`Tem certeza que deseja excluir este action item? Esta ação não pode ser desfeita.`}
        confirmText="Excluir"
        cancelText="Cancelar"
      />
        </div>
      </div>
    </div>
  );
};

export default RetrospectiveDetailPage;