import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { useUrlState } from './lib/urlstate';
import { QuestionList } from './components/QuestionList';
import { Question } from './components/Question';

function App() {
  const [urlState, updateUrlState, revertUrlState] = useUrlState();
  return (
    <div className="App">
      { urlState.id !== null && urlState.id !== undefined ? (
        <Question id={urlState.id} onClose={() => updateUrlState({ id: null })} />
      ) : (
        <QuestionList onCreate={() => updateUrlState({ id: '' })} onEdit={id => updateUrlState({ id })} />
      ) }
    </div>
  );
}

export default App;
