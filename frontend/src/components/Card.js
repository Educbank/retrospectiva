import React from 'react';

const Card = ({
  children,
  className = '',
  flat = false,
  hover = false,
  padding = 'p-6',
  ...props
}) => {
  const baseClasses = flat ? 'card-flat' : 'card';
  const hoverClasses = hover ? 'hover:shadow-md transition-shadow duration-200' : '';
  
  return (
    <div 
      className={`${baseClasses} ${hoverClasses} ${padding} ${className}`}
      {...props}
    >
      {children}
    </div>
  );
};

const CardHeader = ({ children, className = '', ...props }) => (
  <div className={`card-header ${className}`} {...props}>
    {children}
  </div>
);

const CardTitle = ({ children, className = '', ...props }) => (
  <h3 className={`card-title ${className}`} {...props}>
    {children}
  </h3>
);

const CardSubtitle = ({ children, className = '', ...props }) => (
  <p className={`card-subtitle ${className}`} {...props}>
    {children}
  </p>
);

const CardContent = ({ children, className = '', ...props }) => (
  <div className={className} {...props}>
    {children}
  </div>
);

const CardFooter = ({ children, className = '', ...props }) => (
  <div className={`flex items-center justify-between ${className}`} {...props}>
    {children}
  </div>
);

Card.Header = CardHeader;
Card.Title = CardTitle;
Card.Subtitle = CardSubtitle;
Card.Content = CardContent;
Card.Footer = CardFooter;

export default Card;
