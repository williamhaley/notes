import type { Note, UpdateNoteRequest } from './stores/notes';
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

  return res.status;
}

export async function deleteNote(noteId: string) {
  const res = await request(`/api/notes/${noteId}`, {
    method: 'DELETE'
  });

  if (!res.ok) {
    throw new Error('oh no');
  }
}

export async function updateNote(note: UpdateNoteRequest): Promise<Note> {
  const res = await request('/api/notes', {
    method: 'PUT',
    body: JSON.stringify(note)
  });

  if (!res.ok) {
    throw new Error('oh no');
  }

  return await res.json();
}

export async function createNote(): Promise<Note> {
  const res = await request('/api/notes', {
    method: 'POST',
    body: JSON.stringify({
      body: ''
    })
  });

  if (!res.ok) {
    throw new Error('oh no');
  }

  return await res.json();
}

export async function getNotes() {
  const res = await request('/api/notes');

  if (!res.ok) {
    return [];
  }

  const json = (await res.json()) as { results: Note[] };

  return json.results || [];
}
