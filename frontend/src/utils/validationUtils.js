// Validation utility functions

export const validationRules = {
  email: {
    required: true,
    email: true,
    requiredMessage: 'Email é obrigatório',
    emailMessage: 'Email inválido',
  },
  password: {
    required: true,
    minLength: 6,
    requiredMessage: 'Senha é obrigatória',
    minLengthMessage: 'Mínimo 6 caracteres',
  },
  name: {
    required: true,
    minLength: 2,
    requiredMessage: 'Nome é obrigatório',
    minLengthMessage: 'Mínimo 2 caracteres',
  },
  confirmPassword: {
    required: true,
    requiredMessage: 'Confirmação de senha é obrigatória',
  },
};

export const validatePasswordMatch = (password, confirmPassword) => {
  if (password !== confirmPassword) {
    return 'As senhas não coincidem';
  }
  return null;
};

export const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

export const validateRequired = (value, fieldName) => {
  if (!value || value.toString().trim() === '') {
    return `${fieldName} é obrigatório`;
  }
  return null;
};
