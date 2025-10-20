// Date utility functions

export const getRelativeTime = (dateString) => {
  const now = new Date();
  const date = new Date(dateString);
  const diffInMs = now - date;
  const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
  const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

  if (diffInMinutes < 1) {
    return 'Agora mesmo';
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes} minuto${diffInMinutes > 1 ? 's' : ''} atrás`;
  } else if (diffInHours < 24) {
    return `${diffInHours} hora${diffInHours > 1 ? 's' : ''} atrás`;
  } else if (diffInDays < 7) {
    return `${diffInDays} dia${diffInDays > 1 ? 's' : ''} atrás`;
  } else {
    return date.toLocaleDateString('pt-BR');
  }
};

export const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString.split('T')[0] + 'T00:00:00');
  return date.toLocaleDateString('pt-BR');
};

export const isOverdue = (dueDate, status) => {
  if (!dueDate || status === 'done') return false;
  const today = new Date();
  // Parse the date string to avoid timezone issues
  const due = new Date(dueDate.split('T')[0] + 'T00:00:00');
  today.setHours(0, 0, 0, 0);
  due.setHours(0, 0, 0, 0);
  return due < today;
};

export const getDueDateColor = (dueDate, status) => {
  if (!dueDate) return 'text-gray-500';
  if (status === 'done') return 'text-green-600';
  if (isOverdue(dueDate, status)) return 'text-red-600';
  return 'text-blue-600';
};

export const extractCompletionDate = (description) => {
  if (!description) return null;
  
  const dateMatch = description.match(/## ✅ Action Item Concluído - (.*?)\n/);
  return dateMatch ? dateMatch[1].trim() : null;
};
