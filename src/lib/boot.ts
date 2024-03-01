import { writable } from 'svelte/store';
import * as api from './api';
import { init, type Note } from './stores/notes';
import { setToken } from './stores/auth';

export default async function boot() {
  const existingToken = localStorage.getItem('token');

  let notes: Note[] = [];

  if (existingToken) {
    const status = await api.checkToken(existingToken);

    // We are unauthenticated or unauthorized. Not just a matter of the server being down.
    if (status === 401 || status === 403) {
      localStorage.removeItem('token');
      return;
    }

    setToken(existingToken);

    // TODO redundant behavior with notes.ts::getNotes
    notes = await api.getNotes();

    init(notes);
  }
}
