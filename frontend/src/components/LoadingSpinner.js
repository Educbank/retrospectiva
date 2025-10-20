import React from 'react';

const LoadingSpinner = ({ 
  size = 'md', 
  color = 'primary', 
  className = '',
  text = null 
}) => {
  const sizeClasses = {
    xs: 'h-4 w-4',
    sm: 'h-5 w-5',
    md: 'h-8 w-8',
    lg: 'h-12 w-12',
    xl: 'h-16 w-16',
    '2xl': 'h-32 w-32',
  };

  const colorClasses = {
    primary: 'border-primary-600',
    white: 'border-white',
    gray: 'border-gray-600',
    purple: 'border-purple-600',
  };

  const spinnerClasses = `animate-spin rounded-full border-b-2 ${sizeClasses[size]} ${colorClasses[color]} ${className}`;

  if (text) {
    return (
      <div className="flex items-center justify-center">
        <div className={spinnerClasses}></div>
        {text && <span className="ml-2 text-sm text-gray-600">{text}</span>}
      </div>
    );
  }

  return <div className={spinnerClasses}></div>;
};

export default LoadingSpinner;
