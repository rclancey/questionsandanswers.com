import { useState, useCallback, useEffect } from 'react';

const getInitialState = () => {
  const state = {};
  const params = new URLSearchParams(document.location.hash.replace(/^#/, ''));
  Array.from(params.entries()).forEach(([key, val]) => {
    if (state[key] !== undefined) {
      if (Array.isArray(state[key])) {
        state[key].push(val);
      } else if (val !== null && val !== undefined) {
        state[key] = [state[key], val];
      }
    } else {
      state[key] = val;
    }
  });
  state.title = document.title;
  return state;
};

const setHashState = (state) => {
  const params = new URLSearchParams();
  Object.entries(state).forEach(([key, val]) => {
    if (key !== 'title') {
      if (Array.isArray(val)) {
        val.forEach(v => params.append(key, v));
      } else if (val !== null && val !== undefined) {
        params.set(key, val);
      }
    }
  });
  const url = new URL(document.location);
  url.hash = `#${params.toString()}`;
  window.history.pushState(state, state.title || document.title, url.toString());
};

export const useUrlState = () => {
  window.getInitialState = getInitialState;
  const [state, setState] = useState(getInitialState());
  const updateState = useCallback(update => {
    const newState = Object.assign({}, state, update);
    setState(newState);
    setHashState(newState);
  });
  const revertState = useCallback(() => window.history.popState(), []);
  useEffect(() => {
    const onHashChange = () => {
      setState(getInitialState());
    };
    window.addEventListener('hashchange', onHashChange);
    return () => {
      window.removeEventListener('hashchange', onHashChange);
    };
  }, []);
  return [state, updateState, revertState];
};
