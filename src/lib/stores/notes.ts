import { writable } from 'svelte/store';
import * as api from '../api';

export type Note = {
  id: string;
  body?: string;
  timeCreated: number;
  timeLastModified: number;
};

export type UpdateNoteRequest = {
  id: string;
  body: string;
};

export type NoteId = Note['id'];

export type Notes = { [key: NoteId]: Note };

async function createStore() {
  const notesStore = writable({ notes: {} as Notes, selectedNoteId: null as NoteId | null, sortedNoteIds: [] as string[] });

  return {
    init(notes: Note[]) {
      const byKey = notes.reduce((acc, note) => ({ ...acc, [note.id]: note }), {});
      const sortedNoteIds = notes.map(note => note.id);

      notesStore.set({ notes: byKey, selectedNoteId: null, sortedNoteIds });
    },
    async createNote() {
      const note = await api.createNote();

      notesStore.update(store => {
        store.notes = { ...store.notes, [note.id]: note };
        store.selectedNoteId = note.id;
        return store;
      });

      return note;
    },
    async deleteNote(noteId: string) {
      await api.deleteNote(noteId);

      unsetNote(noteId);
    },
    unsetNote(noteId: string) {
      notesStore.update(store => {
        const { [noteId]: _, ...notes } = store.notes;
        store.notes = notes;
        store.selectedNoteId = null;
        store.sortedNoteIds = [...store.sortedNoteIds.filter(id => id !== noteId)];
        return store;
      });
    },
    async getNotes() {
      const notes = await api.getNotes();

      const byKey = notes.reduce((acc, note) => ({ ...acc, [note.id]: note }), {});
      const sortedNoteIds = notes.map(note => note.id);

      notesStore.update(store => {
        store.notes = byKey;
        store.sortedNoteIds = sortedNoteIds;
        return store;
      });
    },
    selectNote(id: NoteId | null) {
      notesStore.update(store => {
        store.selectedNoteId = id;

        return store;
      });
    },
    async updateNote(note: UpdateNoteRequest) {
      const updatedNote = await api.updateNote(note);

      setNote(updatedNote);
    },
    setNote(note: Note) {
      notesStore.update(store => {
        store.notes = {
          ...store.notes,
          [note.id]: note
        };
        store.sortedNoteIds = [note.id, ...store.sortedNoteIds.filter(id => id !== note.id)];

        return store;
      });
    },
    subscribe: notesStore.subscribe
  };
}

const notes = await createStore();

let conn: WebSocket;

type BroadcastEvent = {
  data: string;
};

type BroadcastData = {
  event: string;
  payload: unknown;
};

function observeChanges() {
  conn = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}:${location.port}/socket`);

  conn.onclose = function (evt) {
    console.log('closed', evt);
  };

  conn.onmessage = function (event: BroadcastEvent) {
    const data = JSON.parse(event.data) as BroadcastData;

    switch (data.event) {
      case 'update':
        setNote(data.payload as Note);
        break;
      case 'delete':
        unsetNote(data.payload as string);
        break;
      default:
        break;
    }
  };
}

const { createNote, deleteNote, getNotes, init, selectNote, setNote, subscribe, unsetNote, updateNote } = notes;

export { createNote, deleteNote, getNotes, init, observeChanges, selectNote, setNote, subscribe, unsetNote, updateNote };
