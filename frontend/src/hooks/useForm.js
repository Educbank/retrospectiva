import { useState, useCallback } from 'react';

const useForm = (initialValues = {}, validationRules = {}) => {
  const [formData, setFormData] = useState(initialValues);
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const handleChange = useCallback((e) => {
    const { name, value, type, checked } = e.target;
    const newValue = type === 'checkbox' ? checked : value;
    
    setFormData(prev => ({
      ...prev,
      [name]: newValue,
    }));

    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: '',
      }));
    }
  }, [errors]);

  const setValue = useCallback((name, value) => {
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  }, []);

  const setError = useCallback((name, error) => {
    setErrors(prev => ({
      ...prev,
      [name]: error,
    }));
  }, []);

  const clearErrors = useCallback(() => {
    setErrors({});
  }, []);

  const reset = useCallback(() => {
    setFormData(initialValues);
    setErrors({});
    setLoading(false);
  }, [initialValues]);

  const validate = useCallback(() => {
    const newErrors = {};
    let isValid = true;

    Object.keys(validationRules).forEach(field => {
      const rules = validationRules[field];
      const value = formData[field];

      // Required validation
      if (rules.required && (!value || value.toString().trim() === '')) {
        newErrors[field] = rules.requiredMessage || `${field} é obrigatório`;
        isValid = false;
      }

      // Email validation
      if (rules.email && value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        newErrors[field] = rules.emailMessage || 'Email inválido';
        isValid = false;
      }

      // Min length validation
      if (rules.minLength && value && value.length < rules.minLength) {
        newErrors[field] = rules.minLengthMessage || `Mínimo ${rules.minLength} caracteres`;
        isValid = false;
      }

      // Max length validation
      if (rules.maxLength && value && value.length > rules.maxLength) {
        newErrors[field] = rules.maxLengthMessage || `Máximo ${rules.maxLength} caracteres`;
        isValid = false;
      }

      // Custom validation
      if (rules.custom && value) {
        const customError = rules.custom(value, formData);
        if (customError) {
          newErrors[field] = customError;
          isValid = false;
        }
      }
    });

    setErrors(newErrors);
    return isValid;
  }, [formData, validationRules]);

  const handleSubmit = useCallback(async (e, onSubmit) => {
    e.preventDefault();
    
    if (!validate()) {
      return false;
    }

    setLoading(true);
    try {
      const result = await onSubmit(formData);
      return result;
    } catch (error) {
      console.error('Form submission error:', error);
      return false;
    } finally {
      setLoading(false);
    }
  }, [formData, validate]);

  return {
    formData,
    errors,
    loading,
    handleChange,
    setValue,
    setError,
    clearErrors,
    reset,
    validate,
    handleSubmit,
    setLoading,
  };
};

export default useForm;
