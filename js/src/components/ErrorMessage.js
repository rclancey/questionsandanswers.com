import React from 'react';

export const ErrorMessage = ({ err, onClose }) => (
  <div className="errorMessage">
    <p>{`${err} `}<span onClick={onClose}>x</span></p>
  </div>
);
