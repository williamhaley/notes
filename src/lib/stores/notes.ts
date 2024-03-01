import { writable } from 'svelte/store';
import * as api from '../api';

export type Note = {
  id: string;
  title: string;
  body?: string;
};

async function createStore() {
  const notesStore = writable({ notes: [] as Note[], selectedNoteIndex: -1 });

  return {
    init(notes: Note[]) {
      notesStore.set({ notes, selectedNoteIndex: -1 });
    },
    async createNote() {
      const note = await api.createNote();

      notesStore.update(store => {
        store.notes = [note, ...store.notes];
        store.selectedNoteIndex = 0;
        return store;
      });

      return note;
    },
    async deleteNote(noteId: string) {
      await api.deleteNote(noteId);

      notesStore.update(store => {
        store.notes = store.notes.filter((_, index) => index !== store.selectedNoteIndex);
        store.selectedNoteIndex = -1;
        return store;
      });
    },
    async getNotes() {
      const notes = await api.getNotes();

      notesStore.update(store => {
        store.notes = notes;
        store.selectedNoteIndex = notes.length ? 0 : -1;
        return store;
      });
    },
    selectNote(index: number) {
      notesStore.update(store => {
        store.selectedNoteIndex = index;

        return store;
      });
    },
    async updateNote(note: Note) {
      await api.updateNote(note);

      notesStore.update(store => {
        store.notes = store.notes.with(store.selectedNoteIndex, note);

        return store;
      });
    },
    subscribe: notesStore.subscribe
  };
}

const notes = await createStore();

const { createNote, deleteNote, getNotes, init, selectNote, subscribe, updateNote } = notes;

export { createNote, deleteNote, getNotes, init, selectNote, subscribe, updateNote };
