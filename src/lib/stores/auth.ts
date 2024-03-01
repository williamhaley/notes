import { derived, get, writable } from 'svelte/store';
import { getNotes } from './notes';

const auth = writable('');

export const token = get(auth);

export const getToken = () => {
  return get(auth);
};

export const setToken = (token: string) => {
  auth.set(token);
};

export async function logIn(username: string, password: string) {
  const response = await fetch('/api/login', {
    method: 'POST',
    headers: {
      'content-type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });

  if (response.ok) {
    const token = await response.text();
    setToken(token);

    // TODO better encapsulate this elsewhere
    localStorage.setItem('token', token);

    await getNotes();
  }
}

export const isAuthenticated = derived(auth, store => Boolean(store));
