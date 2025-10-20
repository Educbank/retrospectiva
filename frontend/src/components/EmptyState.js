import React from 'react';

const EmptyState = ({
  icon: Icon,
  title,
  description,
  action,
  className = '',
}) => {
  return (
    <div className={`text-center py-12 ${className}`}>
      {Icon && (
        <Icon className="mx-auto h-12 w-12 text-gray-400" />
      )}
      
      {title && (
        <h3 className="mt-2 text-sm font-medium text-gray-900">
          {title}
        </h3>
      )}
      
      {description && (
        <p className="mt-1 text-sm text-gray-500">
          {description}
        </p>
      )}
      
      {action && (
        <div className="mt-6">
          {action}
        </div>
      )}
    </div>
  );
};

export default EmptyState;
