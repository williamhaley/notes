<script lang="ts">
  import Editor from './Editor.svelte';
  import MediaQuery from './MediaQuery.svelte';
  import Notes from './Notes.svelte';
  import { subscribe, type Note, updateNote } from './stores/notes';

  let note: Note | undefined;

  subscribe(state => {
    note = state.notes[state.selectedNoteIndex];
  });

  function onChange(event: CustomEvent<string>) {
    const draft = { ...note, body: event.detail } as Note;
    updateNote(draft);
  }
</script>

<MediaQuery query="(min-width: 680px)" let:matches>
  {#if matches}
    <div class="split-view">
      <div class="list">
        <Notes />
      </div>
      <div class="details">
        {#if note}
          <Editor body={note.body} id={note.id} on:change={onChange} />
        {/if}
      </div>
    </div>
  {:else}
    {#if note}
      <Editor body={note.body} id={note.id} on:change={onChange} />
    {:else}
      <Notes />
    {/if}
  {/if}
</MediaQuery>



<style>
  .split-view {
    display: flex;
    flex-direction: row;
  }

  .details {
    flex: 1;
    height: 100%;
  }
</style>