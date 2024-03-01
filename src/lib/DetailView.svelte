<script lang="ts">
  import Header from './Header.svelte';
  import { deleteNote, selectNote, type Note, subscribe } from './stores/notes';
  import MediaQuery from './MediaQuery.svelte';
  import Editor from './Editor.svelte';
  import InfoPopover from './InfoPopover.svelte';

  let note: Note;

  subscribe(state => {
    note = state.notes[state.selectedNoteId!];
  });
</script>

<div class="editor-container">
  <Header>
    <div>
      <MediaQuery query="(min-width: 680px)" let:matches>
        {#if !matches}
          <button class="btn btn-primary" type="button" on:click={() => selectNote(null)}>Back</button>
        {/if}
      </MediaQuery>
    </div>
    <div>
      {#if note}
        <InfoPopover {note} />
        <button class="btn btn-sm btn-danger" on:click={() => confirm('Are you sure?') && deleteNote(note.id)}>Delete</button>
      {/if}
    </div>
  </Header>
  {#if note}
    <Editor {note} />
  {/if}
</div>

<style>
  .editor-container {
    height: 100%;
    width: 100%;
    display: grid;
    grid-template-rows: auto 1fr;
  }
</style>
