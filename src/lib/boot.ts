import { writable } from 'svelte/store';
import * as api from './api';
import { init, type Note } from './stores/notes';
import { setToken } from './stores/auth';

export default async function boot() {
  const existingToken = localStorage.getItem('token');

  let notes: Note[] = [];

  if (existingToken) {
    const ok = await api.checkToken(existingToken);

    if (!ok) {
      localStorage.removeItem('token');
    }

    setToken(existingToken);

    // TODO redundant behavior with notes.ts::getNotes
    notes = await api.getNotes();

    init(notes);
  }
}
