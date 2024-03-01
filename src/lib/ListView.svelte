<script lang="ts">
  import Header from './Header.svelte';
  import NoteItem from './NoteItem.svelte';
  import { subscribe, type Notes, selectNote, createNote, type NoteId } from './stores/notes';
  import { toggleMenu } from './stores/ui';

  let notes: Notes = {};
  let sortedNoteIds: string[] = [];
  let selectedNoteId: NoteId | null = null;

  export let menuId: string;

  subscribe(state => {
    notes = state.notes;
    sortedNoteIds = state.sortedNoteIds;
    selectedNoteId = state.selectedNoteId;
  });
</script>

<div class="list">
  <Header>
    <div><button class="btn btn-primary" type="button" aria-controls={menuId} on:click={toggleMenu}>Menu</button></div>
    <div class="d-flex align-items-center fs-4">Notes</div>
    <div><button type="button" class="btn btn-outline-primary" on:click={createNote}>New</button></div>
  </Header>
  <ul class="list-group list-group-flush">
    {#each sortedNoteIds as noteId}
      <a href="/#" class="list-group-item list-group-item-action" class:active={noteId === selectedNoteId} on:click={() => selectNote(noteId)}>
        <NoteItem note={notes[noteId]} />
      </a>
    {/each}
  </ul>
</div>

<style>
  .list {
    height: 100vh;
    overflow: hidden;
  }

  .list-group {
    height: 100%;
    overflow-y: auto;
  }
</style>
