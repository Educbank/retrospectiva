import React from 'react';
import { useQuery } from 'react-query';
import { FileText, Users, Clock, CheckCircle } from 'lucide-react';
import { templatesAPI } from '../services/api';

const TemplatesPage = () => {
  const { data: templates, isLoading } = useQuery('templates', templatesAPI.getTemplates);

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-8 bg-gray-200 rounded w-1/4"></div>
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="card">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-gray-200 rounded w-full mb-2"></div>
              <div className="h-3 bg-gray-200 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  const getTemplateIcon = (templateId) => {
    switch (templateId) {
      case 'start_stop_continue': return '🔄';
      case '4ls': return '📚';
      case 'mad_sad_glad': return '😊';
      case 'sailboat': return '⛵';
      default: return '📋';
    }
  };

  const getTemplateColor = (templateId) => {
    switch (templateId) {
      case 'start_stop_continue': return 'bg-green-100 text-green-800';
      case '4ls': return 'bg-blue-100 text-blue-800';
      case 'mad_sad_glad': return 'bg-purple-100 text-purple-800';
      case 'sailboat': return 'bg-orange-100 text-orange-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Templates de Retrospectiva
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Escolha um template para sua próxima retrospectiva
          </p>
        </div>
      </div>

      {/* Templates Grid */}
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {templates?.data?.map((template) => (
          <div key={template.id} className="card hover:shadow-md transition-shadow duration-200">
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0">
                <div className="h-12 w-12 bg-gray-100 rounded-lg flex items-center justify-center text-2xl">
                  {getTemplateIcon(template.id)}
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="text-lg font-medium text-gray-900">
                  {template.name}
                </h3>
                <p className="text-sm text-gray-500 mt-1">
                  {template.description}
                </p>
                <div className="mt-3">
                  <span className={`badge ${getTemplateColor(template.id)}`}>
                    {template.categories.length} categorias
                  </span>
                </div>
              </div>
            </div>
            
            <div className="mt-4">
              <h4 className="text-sm font-medium text-gray-900 mb-2">Categorias:</h4>
              <div className="space-y-1">
                {template.categories.map((category) => (
                  <div key={category.id} className="flex items-center space-x-2">
                    <div 
                      className="h-3 w-3 rounded-full"
                      style={{ backgroundColor: category.color }}
                    ></div>
                    <span className="text-xs text-gray-600">{category.name}</span>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4 text-xs text-gray-500">
                  <div className="flex items-center">
                    <Users className="h-3 w-3 mr-1" />
                    Qualquer time
                  </div>
                  <div className="flex items-center">
                    <Clock className="h-3 w-3 mr-1" />
                    30-60 min
                  </div>
                </div>
                <button className="btn btn-primary btn-sm">
                  Usar Template
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Template Details Modal would go here */}
      <div className="mt-8 p-6 bg-gray-50 rounded-lg">
        <h3 className="text-lg font-medium text-gray-900 mb-2">
          Como escolher o template certo?
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600">
          <div>
            <h4 className="font-medium text-gray-900 mb-1">Start, Stop, Continue</h4>
            <p>Ideal para times que querem melhorar processos existentes e identificar mudanças específicas.</p>
          </div>
          <div>
            <h4 className="font-medium text-gray-900 mb-1">4Ls (Liked, Learned, Lacked, Longed for)</h4>
            <p>Perfeito para refletir sobre experiências e aprendizados de forma estruturada.</p>
          </div>
          <div>
            <h4 className="font-medium text-gray-900 mb-1">Mad, Sad, Glad</h4>
            <p>Ótimo para explorar aspectos emocionais e criar um ambiente seguro para feedback.</p>
          </div>
          <div>
            <h4 className="font-medium text-gray-900 mb-1">Sailboat</h4>
            <p>Ideal para visualizar o progresso em direção aos objetivos e identificar obstáculos.</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TemplatesPage;
