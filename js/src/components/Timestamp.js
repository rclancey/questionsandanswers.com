import React from 'react';
import moment from 'moment';

export const Timestamp = ({ children }) => (
  <div className="timestamp">
    {moment(children).format('MMMM Do YYYY, h:mm:ss a')}
  </div>
);
