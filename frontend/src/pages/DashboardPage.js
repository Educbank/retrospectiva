import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from 'react-query';
import { 
  Users, 
  MessageSquare, 
  Plus, 
  TrendingUp, 
  Calendar,
  Activity,
  Clock,
  Filter,
  CheckSquare,
  AlertCircle
} from 'lucide-react';
import { usersAPI, retrospectivesAPI } from '../services/api';
import Card from '../components/Card';
import EmptyState from '../components/EmptyState';
import LoadingSpinner from '../components/LoadingSpinner';
import { 
  getRelativeTime, 
  formatDate, 
  isOverdue, 
  getDueDateColor, 
  extractCompletionDate,
  getActionItemStatusColor,
  getActionItemStatusText
} from '../utils';

const DashboardPage = () => {
  const { data: user } = useQuery('userProfile', usersAPI.getProfile);
  const { data: retrospectives, isLoading: retrospectivesLoading } = useQuery('userRetrospectives', retrospectivesAPI.getRetrospectives);
  
  // States for Action Items filters
  const [statusFilter, setStatusFilter] = useState('all');
  const [hideCompleted, setHideCompleted] = useState(false);

  // Calculate total action items from retrospectives (considering visibility filter)
  const totalActionItems = useMemo(() => {
    if (!retrospectives?.data || !user?.data) return 0;
    
    return retrospectives.data.reduce((total, retro) => {
      // Filter retrospectives: show "planejada" only to creator, others to everyone
      const shouldShowRetro = retro.status !== 'planejada' || retro.created_by === user.data.id;
      
      if (shouldShowRetro) {
        return total + (retro.action_items?.length || 0);
      }
      return total;
    }, 0);
  }, [retrospectives?.data, user?.data]);

  // Process and filter Action Items
  const allActionItems = useMemo(() => {
    if (!retrospectives?.data || !user?.data) return [];
    
    const items = [];
    retrospectives.data.forEach(retro => {
      // Filter retrospectives: show "planejada" only to creator, others to everyone
      const shouldShowRetro = retro.status !== 'planejada' || retro.created_by === user.data.id;
      
      if (shouldShowRetro && retro.action_items) {
        retro.action_items.forEach(item => {
          items.push({
            ...item,
            retrospective: {
              id: retro.id,
              title: retro.title
            }
          });
        });
      }
    });
    
    // Sort by creation date (newest first)
    return items.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  }, [retrospectives?.data, user?.data]);

  // Filter Action Items based on status and hide completed
  const filteredActionItems = useMemo(() => {
    let filtered = allActionItems;
    
    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter(item => item.status === statusFilter);
    }
    
    // Hide completed if checkbox is checked
    if (hideCompleted) {
      filtered = filtered.filter(item => item.status !== 'done');
    }
    
    return filtered;
  }, [allActionItems, statusFilter, hideCompleted]);

  // Calculate visible retrospectives count (considering visibility filter)
  const visibleRetrospectives = useMemo(() => {
    if (!retrospectives?.data || !user?.data) return [];
    
    return retrospectives.data.filter(retro => {
      // Filter retrospectives: show "planejada" only to creator, others to everyone
      return retro.status !== 'planejada' || retro.created_by === user.data.id;
    });
  }, [retrospectives?.data, user?.data]);

  const stats = [
    {
      name: 'Retrospectivas',
      value: visibleRetrospectives.length,
      icon: MessageSquare,
      iconBg: 'bg-gray-100',
      iconColor: 'text-gray-600',
      link: '/retrospectives',
    },
    {
      name: 'Em Andamento',
      value: visibleRetrospectives.filter(r => r.status === 'active').length,
      icon: Activity,
      iconBg: 'bg-blue-50',
      iconColor: 'text-blue-600',
      link: '/retrospectives',
    },
    {
      name: 'Concluídas',
      value: visibleRetrospectives.filter(r => r.status === 'closed').length,
      icon: TrendingUp,
      iconBg: 'bg-green-50',
      iconColor: 'text-green-600',
      link: '/retrospectives',
    },
    {
      name: 'Action Items',
      value: totalActionItems,
      icon: Clock,
      iconBg: 'bg-purple-50',
      iconColor: 'text-purple-600',
      link: '/action-items',
    },
  ];


  const recentActivities = visibleRetrospectives?.slice(0, 3).map((retrospective) => ({
    id: retrospective.id,
    type: 'retrospective',
    title: retrospective.title,
    team: retrospective.template,
    time: getRelativeTime(retrospective.created_at),
    status: retrospective.status,
  })) || [
    {
      id: 1,
      type: 'retrospective',
      title: 'Sprint Review - Q1 2024',
      team: 'sailboat',
      time: '2 horas atrás',
      status: 'em_andamento',
    },
    {
      id: 2,
      type: 'retrospective',
      title: 'Retrospectiva Mensal',
      team: 'start_stop_continue',
      time: '1 dia atrás',
      status: 'concluido',
    },
    {
      id: 3,
      type: 'retrospective',
      title: 'Review de Projeto',
      team: '4ls',
      time: '2 dias atrás',
      status: 'concluido',
    },
  ];

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <Card className="bg-white rounded-xl p-8 shadow-sm border border-gray-100">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900 mb-1">
            Olá, {user?.data?.name?.split(' ')[0]}!
          </h1>
          <p className="text-gray-600">Resumo das suas atividades</p>
        </div>
      </Card>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Link
              key={stat.name}
              to={stat.link}
              className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow duration-200 border border-gray-100"
            >
              <div className="flex items-center space-x-4">
                <div className={`flex items-center justify-center w-10 h-10 rounded-lg ${stat.iconBg}`}>
                  <Icon className={`h-5 w-5 ${stat.iconColor}`} />
                </div>
                <div>
                  <div className="text-2xl font-semibold text-gray-900">
                    {stat.value}
                  </div>
                  <div className="text-sm text-gray-600">
                    {stat.name}
                  </div>
                </div>
              </div>
            </Link>
          );
        })}
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
        {/* Recent Retrospectives */}
        <Card>
          <Card.Header>
            <Card.Title>Últimas Retrospectivas</Card.Title>
          </Card.Header>
          <Card.Content>
            <div className="space-y-4">
              {retrospectivesLoading ? (
                <div className="animate-pulse space-y-3">
                  {[...Array(3)].map((_, i) => (
                    <div key={i} className="flex items-center space-x-3">
                      <div className="h-10 w-10 bg-gray-200 rounded-lg"></div>
                      <div className="flex-1 space-y-2">
                        <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                        <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : retrospectives?.data?.length > 0 ? (
                retrospectives.data.slice(0, 5).map((retrospective) => (
                  <div key={retrospective.id} className="flex items-center space-x-3 p-3 hover:bg-gray-50 rounded-lg transition-colors">
                    <div className="flex-shrink-0">
                      <div className={`h-8 w-8 rounded-lg flex items-center justify-center ${
                        retrospective.status === 'active' ? 'bg-blue-100' :
                        retrospective.status === 'closed' ? 'bg-green-100' : 'bg-gray-100'
                      }`}>
                        <MessageSquare className={`h-4 w-4 ${
                          retrospective.status === 'active' ? 'text-blue-600' :
                          retrospective.status === 'closed' ? 'text-green-600' : 'text-gray-600'
                        }`} />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {retrospective.title}
                      </p>
                      <p className="text-sm text-gray-500 truncate">
                        {retrospective.template.replace('_', ' ')} • {getRelativeTime(retrospective.created_at)}
                      </p>
                    </div>
                    <div className="flex-shrink-0">
                      <Link
                        to={`/retrospectives/${retrospective.id}`}
                        className="text-purple-600 hover:text-purple-700 text-sm font-medium"
                      >
                        Ver
                      </Link>
                    </div>
                  </div>
                ))
              ) : (
                <EmptyState
                  icon={MessageSquare}
                  title="Nenhuma retrospectiva encontrada"
                  description="Comece criando sua primeira retrospectiva."
                  action={
                    <Link to="/retrospectives/new" className="btn btn-primary">
                      <Plus className="h-4 w-4 mr-2" />
                      Nova Retrospectiva
                    </Link>
                  }
                />
              )}
            </div>
          </Card.Content>
        </Card>

        {/* Action Items */}
        <Card>
          <Card.Header>
            <div className="flex items-center justify-between">
              <Card.Title>Action Items</Card.Title>
              <div className="flex items-center space-x-3">
                <div className="flex items-center space-x-2">
                  <Filter className="h-4 w-4 text-gray-400" />
                  <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                    className="text-sm border border-gray-300 rounded-md px-2 py-1 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="all">Todos</option>
                    <option value="todo">A fazer</option>
                    <option value="in_progress">Em andamento</option>
                    <option value="done">Concluídos</option>
                  </select>
                </div>
                <label className="flex items-center space-x-2 text-sm text-gray-600">
                  <input
                    type="checkbox"
                    checked={hideCompleted}
                    onChange={(e) => setHideCompleted(e.target.checked)}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span>Ocultar concluídos</span>
                </label>
              </div>
            </div>
          </Card.Header>
          <Card.Content>
            <div className="space-y-3">
              {retrospectivesLoading ? (
                <div className="animate-pulse space-y-3">
                  {[...Array(3)].map((_, i) => (
                    <div key={i} className="flex items-center space-x-3 p-3">
                      <div className="h-8 w-8 bg-gray-200 rounded-lg"></div>
                      <div className="flex-1 space-y-2">
                        <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                        <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : filteredActionItems.length > 0 ? (
                filteredActionItems.slice(0, 5).map((actionItem) => (
                  <div key={actionItem.id} className="flex items-center space-x-3 p-3 hover:bg-gray-50 rounded-lg transition-colors">
                    <div className="flex-shrink-0">
                      <div className="h-8 w-8 rounded-lg flex items-center justify-center bg-purple-100">
                        <CheckSquare className="h-4 w-4 text-purple-600" />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {actionItem.title}
                      </p>
                      <div className="flex items-center space-x-2 mt-1">
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getActionItemStatusColor(actionItem.status)}`}>
                          {getActionItemStatusText(actionItem.status)}
                        </span>
                        {isOverdue(actionItem.due_date, actionItem.status) && (
                          <span className="badge badge-danger">
                            <AlertCircle className="h-3 w-3 inline mr-1" />
                            Atrasado
                          </span>
                        )}
                        <span className="text-green-600 text-sm">
                          {actionItem.completed_at ? formatDate(actionItem.completed_at) : extractCompletionDate(actionItem.description)}
                        </span>
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
                        {actionItem.status === 'done' && (actionItem.completed_at || extractCompletionDate(actionItem.description)) && (
                          <div className="flex items-center space-x-1">
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="flex-shrink-0">
                      <Link
                        to="/action-items"
                        className="text-purple-600 hover:text-purple-700 text-sm font-medium"
                      >
                        Ver
                      </Link>
                    </div>
                  </div>
                ))
              ) : (
                <EmptyState
                  icon={CheckSquare}
                  title="Nenhum action item encontrado"
                  description="Nenhum action item com este status."
                />
              )}
            </div>
          </Card.Content>
        </Card>
      </div>
    </div>
  );
};

export default DashboardPage;
