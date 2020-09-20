import React, { useState, useEffect } from 'react';

export const Loading = () => {
  const [dots, setDots] = useState(1);
  useEffect(() => {
    const interval = setInterval(() => {
      setDots(orig => (orig % 5) + 1);
    }, 500);
    return () => {
      clearInterval(interval);
    };
  }, []);
  return (<div className="loading">{'.'.repeat(dots)}</div>);
};
