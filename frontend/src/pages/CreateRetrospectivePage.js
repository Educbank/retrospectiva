import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { MessageSquare, ArrowLeft, Plus } from 'lucide-react';
import { retrospectivesAPI, templatesAPI } from '../services/api';
import toast from 'react-hot-toast';

const CreateRetrospectivePage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    template: ''
  });

  const { data: templates, isLoading: templatesLoading } = useQuery('templates', templatesAPI.getTemplates);

  const createRetrospectiveMutation = useMutation(
    (data) => retrospectivesAPI.createRetrospective(data),
    {
      onSuccess: (response) => {
        queryClient.invalidateQueries('userRetrospectives');
        toast.success('Retrospectiva criada com sucesso!');
        navigate(`/retrospectives/${response.data.id}`);
      },
      onError: (error) => {
        toast.error('Erro ao criar retrospectiva: ' + (error.response?.data?.error || error.message));
      },
    }
  );

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!formData.title.trim()) {
      toast.error('Título é obrigatório');
      return;
    }
    
    if (!formData.template) {
      toast.error('Template é obrigatório');
      return;
    }

    createRetrospectiveMutation.mutate(formData);
  };

  const getTemplateIcon = (templateId) => {
    switch (templateId) {
      case 'start_stop_continue': return '🔄';
      case '4ls': return '📚';
      case 'mad_sad_glad': return '😊';
      case 'sailboat': return '⛵';
      case 'went_well_to_improve': return '📈';
      default: return '📋';
    }
  };

  const getTemplateDescription = (templateId) => {
    switch (templateId) {
      case 'start_stop_continue': return 'Identifique o que começar, parar e continuar fazendo';
      case '4ls': return 'Reflita sobre o que gostou, aprendeu, faltou e deseja';
      case 'mad_sad_glad': return 'Explore aspectos emocionais do trabalho';
      case 'sailboat': return 'Visualize progresso e obstáculos em direção aos objetivos';
      case 'went_well_to_improve': return 'Avalie o que funcionou bem e o que pode ser melhorado';
      default: return 'Template personalizado';
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-6">
        <button
          onClick={() => navigate('/retrospectives')}
          className="flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Voltar para Retrospectivas
        </button>
        
        <div className="flex items-center space-x-3">
          <div className="h-10 w-10 bg-primary-100 rounded-lg flex items-center justify-center">
            <MessageSquare className="h-5 w-5 text-primary-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Nova Retrospectiva</h1>
            <p className="text-gray-600">Crie uma nova retrospectiva para sua equipe</p>
          </div>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">Informações Básicas</h2>
          
          <div className="space-y-4">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                Título da Retrospectiva *
              </label>
              <input
                type="text"
                id="title"
                name="title"
                value={formData.title}
                onChange={handleChange}
                className="input mt-1"
                placeholder="Ex: Sprint Review - Q1 2024"
                required
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Descrição
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                rows={3}
                className="input mt-1"
                placeholder="Descreva o contexto e objetivos desta retrospectiva..."
              />
            </div>

            <div>
              <label htmlFor="template" className="block text-sm font-medium text-gray-700">
                Template *
              </label>
              {templatesLoading ? (
                <div className="mt-1">
                  <div className="animate-pulse h-10 bg-gray-200 rounded-md"></div>
                </div>
              ) : (
                <select
                  id="template"
                  name="template"
                  value={formData.template}
                  onChange={handleChange}
                  className="input mt-1"
                  required
                >
                  <option value="">Selecione um template...</option>
                  {templates?.data?.map((template) => (
                    <option key={template.id} value={template.id}>
                      {getTemplateIcon(template.id)} {template.name}
                    </option>
                  ))}
                </select>
              )}
              
              {formData.template && (
                <div className="mt-3 p-3 bg-gray-50 rounded-md">
                  <p className="text-sm text-gray-600">
                    <span className="font-medium">{getTemplateIcon(formData.template)} {templates?.data?.find(t => t.id === formData.template)?.name}:</span>
                    <br />
                    {getTemplateDescription(formData.template)}
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="flex items-center justify-end space-x-3">
          <button
            type="button"
            onClick={() => navigate('/retrospectives')}
            className="btn btn-secondary"
          >
            Cancelar
          </button>
          <button
            type="submit"
            disabled={createRetrospectiveMutation.isLoading || !formData.title.trim() || !formData.template}
            className="btn btn-primary"
          >
            {createRetrospectiveMutation.isLoading ? (
              <div className="flex items-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                Criando...
              </div>
            ) : (
              <div className="flex items-center">
                <Plus className="h-4 w-4 mr-2" />
                Criar Retrospectiva
              </div>
            )}
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateRetrospectivePage;