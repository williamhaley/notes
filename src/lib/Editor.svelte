<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { deleteNote } from './stores/notes';

  const dispatch = createEventDispatcher();

  export let id: string;
  export let body: string | undefined;

  let textarea: HTMLTextAreaElement;

  let timer: ReturnType<typeof setTimeout>;

	function onChange() {
    clearTimeout(timer);
		timer = setTimeout(() => {
      dispatch('change', textarea.value);
		}, 800);
	}
</script>

<div class="editor-container">
  <div class="controls p-1">
    <button class="btn btn-sm btn-danger" on:click={() => confirm('Are you sure?') && deleteNote(id)}>Delete</button>
  </div>
  <textarea class="editor" bind:this={textarea} bind:value={body} on:input={onChange}
  ></textarea>
  <div>Footer</div>
</div>

<style>
  .editor-container {
    height: 100%;
    width: 100%;
    display: grid;
    grid-template-rows: auto 1fr auto;
  }

  .controls {
    border-bottom: 1px solid #ddd;
  }

  .editor {
    padding: 1em;
    box-sizing: border-box;
    border: none;
    resize: none;
    outline: none;
  }
</style>
