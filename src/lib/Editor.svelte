<script lang="ts">
  import { afterUpdate, onDestroy, onMount } from 'svelte';
  import EasyMDE from 'easymde';
  import { updateNote, type Note } from './stores/notes';

  export let note: Note;

  let textarea: HTMLTextAreaElement;
  let editor: EasyMDE;
  let previousId: string;

  onMount(() => {
    editor = new EasyMDE({
      element: textarea,
      autofocus: true,
      toolbar: false,
      status: false,
      sideBySideFullscreen: false,
      initialValue: note.body
    });

    editor.codemirror.on('change', function () {
      if (editor.value() !== note.body) {
        updateNote({ id: note.id, body: editor.value() });
      }
    });
  });

  onDestroy(() => {
    editor.cleanup();
    editor.toTextArea();
  });

  afterUpdate(() => {
    if (note.id !== previousId) {
      editor.value(note.body ?? '');
      previousId = note.id;
    }
  });
</script>

<textarea class="editor" autocomplete="on" bind:this={textarea}></textarea>

<style>
  :global(.EasyMDEContainer) {
    overflow-y: auto;
  }

  :global(.EasyMDEContainer .CodeMirror) {
    height: 100%;
  }

  .editor {
    visibility: hidden;
  }
</style>
