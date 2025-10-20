// Status utility functions

export const getActionItemStatusColor = (status) => {
  switch (status) {
    case 'done': return 'bg-green-100 text-green-800';
    case 'in_progress': return 'bg-yellow-100 text-yellow-800';
    case 'todo': return 'bg-gray-100 text-gray-800';
    default: return 'bg-gray-100 text-gray-800';
  }
};

export const getActionItemStatusText = (status) => {
  switch (status) {
    case 'done': return 'Concluído';
    case 'in_progress': return 'Em andamento';
    case 'todo': return 'A fazer';
    default: return status;
  }
};

export const getRetrospectiveStatusColor = (status) => {
  switch (status) {
    case 'active': return 'bg-blue-100 text-blue-800';
    case 'closed': return 'bg-green-100 text-green-800';
    case 'planejada': return 'bg-gray-100 text-gray-800';
    default: return 'bg-gray-100 text-gray-800';
  }
};

export const getRetrospectiveStatusText = (status) => {
  switch (status) {
    case 'active': return 'Em andamento';
    case 'closed': return 'Concluída';
    case 'planejada': return 'Planejada';
    default: return status;
  }
};
