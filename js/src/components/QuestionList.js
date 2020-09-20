import React, { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { fetchJson } from '../lib/api';
import { Timestamp } from './Timestamp';

const questionsUrl = ({ sort, page, count }) => {
  const url = `/api/questions?sort=${sort}&page=${page}&count=${count}`;
  return url;
};

const getQuestions = (url) => fetchJson(url, { method: 'GET' }).catch(err => console.error(err));

export const QuestionList = ({ onCreate, onEdit }) => {
  const ref = useRef(null);
  const [sort, setSort] = useState('-questionDate');
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(5);
  const [questions, setQuestions] = useState([]);
  const [nextPage, setNextPage] = useState(null);
  useEffect(() => {
    const url = questionsUrl({ sort, page, count });
    getQuestions(url)
      .then(resp => {
        setNextPage(resp.meta.nextPage);
        setQuestions(resp.results);
      });
  }, [sort, page, count]);
  const onScroll = useCallback(evt => {
    if (nextPage) {
      const frame = evt.target.getBoundingClientRect();
      const last = evt.target.lastElementChild.getBoundingClientRect();
      if (last.top < frame.bottom + frame.height) {
        setNextPage(null);
        getQuestions(nextPage)
          .then(resp => {
            setQuestions(orig => {
              const seen = new Set(orig.map(item => item.id));
              return orig.concat(resp.results.filter(item => !seen.has(item.id)));
            });
            setNextPage(resp.meta.nextPage);
          });
      }
    }
  }, [nextPage]);
  useEffect(() => {
    if (questions.length > 0 && nextPage && ref.current) {
      onScroll({ target: ref.current });
    }
  }, [questions, nextPage]);
  return (
    <>
      <input type="button" value="Ask a question" onClick={onCreate} />
      <div
        className="questionList"
        onScroll={onScroll}
        ref={node => {
          if (node) {
            ref.current = node;
          }
        }}
      >
        {questions.map(q => <Question key={q.id} onEdit={onEdit} {...q} />)}
      </div>
    </>
  );
};

const Question = ({ id, question, questionDate, answer, answerDate, onEdit }) => (
  <div className="qAndA">
    <div className="question" onClick={() => onEdit(id)}>
      <Timestamp>{questionDate}</Timestamp>
      <p>{question}</p>
    </div>
    { answer ? (
      <div className="answer">
        <Timestamp>{answerDate}</Timestamp>
        <p>{answer}</p>
      </div>
    ) : (
      <div className="answer" />
    ) }
  </div>
);
