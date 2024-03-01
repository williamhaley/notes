<script lang="ts">
  import { format } from 'date-fns';
  import type { Note } from './stores/notes';
  import { afterUpdate, onMount } from 'svelte';
  import { Popover } from 'bootstrap';

  export let note: Note;

  let popoverButton: HTMLAnchorElement;
  let popover: Popover;

  function formatDate(timestamp: number) {
    return format(new Date(timestamp * 1000), 'yyyy-MM-dd p');
  }

  onMount(() => {
    popover = new Popover(popoverButton, {
      trigger: 'focus',
      html: true
    });
  });

  afterUpdate(() => {
    popover.setContent({
      '.popover-body': `Created: ${formatDate(note.timeCreated)}<br>Modified: ${formatDate(note.timeLastModified)}`
    });
  });
</script>

<a href="/#" tabindex={0} bind:this={popoverButton} type="button" class="btn btn-sm btn-info me-1" data-bs-toggle="popover" data-bs-title="Info"
  >Info</a
>
