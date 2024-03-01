import type { Note } from './stores/notes';
import request from './request';

export async function logIn(username: string, password: string) {
  const res = await fetch('/api/login', {
    method: 'POST',
    body: JSON.stringify({ username, password })
  });

  if (res.status !== 200) {
    throw new Error('nope');
  }

  return await res.text();
}

export async function checkToken(token: string) {
  const res = await request('/api/check', {
    headers: {
      token
    }
  });

  return res.ok;
}

export async function deleteNote(noteId: string) {
  const res = await request(`/api/indexes/notes/documents/${noteId}`, {
    method: 'DELETE'
  });

  if (!res.ok) {
    throw new Error('oh no');
  }
}

export async function updateNote(note: Note) {
  const res = await request('/api/indexes/notes/documents', {
    method: 'POST',
    body: JSON.stringify(note)
  });

  if (!res.ok) {
    throw new Error('oh no');
  }
}

export async function createNote() {
  const note = {
    id: Math.random().toString(36).slice(2, 7),
    title: '',
    body: ''
  };

  const res = await request('/api/indexes/notes/documents', {
    method: 'POST',
    body: JSON.stringify(note)
  });

  if (!res.ok) {
    throw new Error('oh no');
  }

  return note;
}

export async function getNotes() {
  const res = await request('/api/indexes/notes/documents');

  if (!res.ok) {
    return [];
  }

  const json = (await res.json()) as { results: Note[] };

  return json.results || [];
}
