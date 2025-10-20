import React, { useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';

const InputField = ({
  label,
  name,
  type = 'text',
  placeholder,
  value,
  onChange,
  error,
  required = false,
  icon: Icon,
  autoComplete,
  className = '',
  disabled = false,
}) => {
  const [showPassword, setShowPassword] = useState(false);
  const isPassword = type === 'password';
  const inputType = isPassword && showPassword ? 'text' : type;

  return (
    <div className="space-y-2">
      {label && (
        <label htmlFor={name} className="block text-sm font-medium text-gray-700">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      <div className="relative">
        {Icon && (
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Icon className="h-5 w-5 text-gray-400" />
          </div>
        )}
        
        <input
          id={name}
          name={name}
          type={inputType}
          autoComplete={autoComplete}
          required={required}
          value={value}
          onChange={onChange}
          disabled={disabled}
          className={`input ${Icon ? 'pl-10' : ''} ${isPassword ? 'pr-10' : ''} ${
            error ? 'input-error' : ''
          } ${className}`}
          placeholder={placeholder}
        />
        
        {isPassword && (
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-3 flex items-center"
            onClick={() => setShowPassword(!showPassword)}
            tabIndex={-1}
          >
            {showPassword ? (
              <EyeOff className="h-5 w-5 text-gray-400 hover:text-gray-600 transition-colors" />
            ) : (
              <Eye className="h-5 w-5 text-gray-400 hover:text-gray-600 transition-colors" />
            )}
          </button>
        )}
      </div>
      
      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}
    </div>
  );
};

export default InputField;
