import React, { useState, useEffect, useCallback } from 'react';
import { fetchJson, HTTPError } from '../lib/api';
import { Timestamp } from './Timestamp';
import { ErrorMessage } from './ErrorMessage';
import { Loading } from './Loading';

export const Question = ({ id, onClose }) => {
  const [myId, setMyId] = useState(id);
  const [editing, setEditing] = useState(!myId);
  const [error, setError] = useState(null);
  const [q, setQ] = useState({});
  useEffect(() => setEditing(orig => orig || !myId), [myId]);
  useEffect(() => {
    setMyId(id);
    if (id) {
      fetchJson(`/api/question/${id}`, { method: 'GET' })
        .then(obj => setQ(obj))
        .catch(err => setError(err));
    }
  }, [id]);
  useEffect(() => {
    if (q.id !== undefined) {
      setMyId(q.id);
      setEditing(false);
    }
  }, [q]);
  if (error) {
    return (<ErrorMessage err={error} onClose={onClose} />);
  }
  if (!myId) {
    return (<CreateQuestion onSave={onClose} onCancel={onClose} />);
  }
  if (!q.id) {
    return (<Loading />);
  }
  if (editing) {
    return (<EditAnswer {...q} onSave={setQ} onCancel={() => setEditing(false)} />);
  }
  return (
    <div className="qAndA">
      <div className="question">
        <h1>Question:</h1>
        <Timestamp>{q.questionDate}</Timestamp>
        <p>{q.question}</p>
      </div>
      { q.answer ? (
        <div className="answer">
          <h1>Answer:</h1>
          <Timestamp>{q.answerDate}</Timestamp>
          <p>{q.answer}</p>
          <input type="button" value="Edit answer" onClick={() => setEditing(true)} />
        </div>
      ) : (
        <div className="noanswer">
          <p>This question has not been answered</p>
          <input type="button" value="Answer it!" onClick={() => setEditing(true)} />
        </div>
      ) }
    </div>
  );
};

const CreateQuestion = ({ onSave, onCancel }) => {
  const [question, setQuestion] = useState('');
  const [error, setError] = useState(null);
  const onCreate = useCallback(() => {
    fetchJson('/api/question', {
      method: 'POST',
      header: {
        'Content-Type': 'text/plain',
      },
      body: question,
    })
      .then(q => onSave(q))
      .catch(err => setError(err));
  }, [question]);
  return (
    <div className="createQuestion">
      { error !== null ? (
        <ErrorMessage err={error} onClose={() => setError(null)} />
      ) : null }
      <textarea
        cols={40}
        rows={10}
        value={question}
        onChange={evt => setQuestion(evt.target.value)}
      />
      <div className="buttons">
        <input type="button" value="Create Question" onClick={onCreate} />
        <input type="button" value="Cancel" onClick={onCancel} />
      </div>
    </div>
  );
};

const EditAnswer = ({ id, question, questionDate, answer, answerDate, onSave, onCancel }) => {
  const [myAnswer, setMyAnswer] = useState(answer);
  const [error, setError] = useState(null);
  useEffect(() => {
    setMyAnswer(answer);
  }, [answer]);
  const onEdit = useCallback(() => {
    fetchJson(`/api/question/${id}`, {
      method: 'PUT',
      header: {
        'Content-Type': 'text/plain',
      },
      body: myAnswer,
    })
      .then(q => onSave(q))
      .catch(err => setError(err));
  }, [id, myAnswer]);
  return (
    <div className="editAnswer">
      <div className="question">
        <Timestamp>{questionDate}</Timestamp>
        <p>{question}</p>
      </div>
      { error !== null ? (
        <ErrorMessage err={error} onClose={() => setError(null)} />
      ) : null }
      <div className="answer">
        <textarea
          cols={40}
          rows={10}
          value={myAnswer}
          onChange={evt => setMyAnswer(evt.target.value)}
        />
      </div>
      <div className="buttons">
        <input type="button" value="Save Answer" onClick={onEdit} />
        <input type="button" value="Cancel" onClick={onCancel} />
      </div>
    </div>
  );
};
