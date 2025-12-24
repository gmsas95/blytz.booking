import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  selected?: boolean;
}

export const Card: React.FC<CardProps> = ({ children, className = '', onClick, selected }) => {
  return (
    <div 
      onClick={onClick}
      className={`
        relative overflow-hidden rounded-lg bg-white p-4 transition-all duration-200
        ${onClick ? 'cursor-pointer hover:border-primary-300' : ''}
        ${selected ? 'ring-2 ring-primary-500 border-primary-500 shadow-md' : 'border border-gray-200 shadow-sm'}
        ${className}
      `}
    >
      {children}
    </div>
  );
};
