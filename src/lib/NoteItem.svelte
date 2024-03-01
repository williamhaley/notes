<script lang="ts">
  import { type Note } from './stores/notes';

  export let note: Note;

  function getTitle({ body }: Note) {
    const [firstLine, ...rest] = (body || '').split('\n').filter(line => line.trim() !== '');

    return {
      title: firstLine?.slice(0, 100) || null,
      preview: rest?.join(' ').slice(0, 100) || null
    };
  }

  $: pieces = getTitle(note);
  $: title = pieces.title;
  $: preview = pieces.preview;
</script>

<div class="wrapper">
  <div class="title">{title || 'New note ...'}</div>
  {#if preview}
    <div class="preview">{preview}</div>
  {/if}
</div>

<style>
  .wrapper {
    width: 100%;
    display: table;
    table-layout: fixed;
  }

  .title {
    font-size: 1.5em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .preview {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
