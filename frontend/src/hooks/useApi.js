import { useQuery, useMutation, useQueryClient } from 'react-query';
import toast from 'react-hot-toast';

const useApi = () => {
  const queryClient = useQueryClient();

  const createQuery = (queryKey, queryFn, options = {}) => {
    return useQuery(queryKey, queryFn, {
      retry: 1,
      refetchOnWindowFocus: false,
      ...options,
    });
  };

  const createMutation = (mutationFn, options = {}) => {
    return useMutation(mutationFn, {
      onError: (error) => {
        const message = error.response?.data?.error || 'Erro inesperado';
        toast.error(message);
      },
      ...options,
    });
  };

  const invalidateQueries = (queryKey) => {
    queryClient.invalidateQueries(queryKey);
  };

  const setQueryData = (queryKey, data) => {
    queryClient.setQueryData(queryKey, data);
  };

  const getQueryData = (queryKey) => {
    return queryClient.getQueryData(queryKey);
  };

  const refetchQueries = (queryKey) => {
    queryClient.refetchQueries(queryKey);
  };

  const createMutationWithDefaults = (mutationFn, options = {}) => {
    return useMutation(mutationFn, {
      onSuccess: (data, variables, context) => {
        if (options.onSuccess) {
          options.onSuccess(data, variables, context);
        }
        
        // Auto-invalidate related queries
        if (options.invalidateQueries) {
          options.invalidateQueries.forEach(queryKey => {
            invalidateQueries(queryKey);
          });
        }
      },
      onError: (error, variables, context) => {
        if (options.onError) {
          options.onError(error, variables, context);
        } else {
          const message = error.response?.data?.error || 'Erro inesperado';
          toast.error(message);
        }
      },
      ...options,
    });
  };

  const createQueryWithDefaults = (queryKey, queryFn, options = {}) => {
    return useQuery(queryKey, queryFn, {
      retry: 1,
      refetchOnWindowFocus: false,
      ...options,
    });
  };

  return {
    useQuery: createQueryWithDefaults,
    useMutation: createMutationWithDefaults,
    invalidateQueries,
    setQueryData,
    getQueryData,
    refetchQueries,
    queryClient,
  };
};

export default useApi;
