<script lang="ts">
  import NoteItem from './NoteItem.svelte';
  import { subscribe, type Note, selectNote } from './stores/notes';

  let notes: Note[] = [];
  let selectedNoteIndex = -1;

  subscribe(state => {
    notes = state.notes;
    selectedNoteIndex = state.selectedNoteIndex;
  });
</script>

<ul class="list-group list-group-flush">
  {#each notes as note, index}
    <a href="#{note.id}" class="list-group-item list-group-item-action" class:active={index === selectedNoteIndex} on:click={() => selectNote(index)}>
      <NoteItem {note} />
    </a>
  {/each}
</ul>