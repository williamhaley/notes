<script lang="ts">
  import { Offcanvas } from 'bootstrap';
  import { onMount } from 'svelte';
  import { subscribe, toggleSidebar } from './stores/ui';

  export let id: string;

  let sidebar: Offcanvas;
  let sidebarElement: HTMLDivElement;

  onMount(() => {
    sidebar = new Offcanvas(sidebarElement);
  });

  subscribe(({ isSidebarOpen }) => {
    if (!sidebar) {
      return;
    }
    
    if (isSidebarOpen) {
      sidebar.show();
    } else {
      sidebar.hide();
    }
  });
</script>

<div bind:this={sidebarElement} class="offcanvas offcanvas-start" tabindex="-1" id={id} aria-labelledby="sidebarLabel">
  <div class="offcanvas-header">
    <h5 id="sidebarLabel">Menu</h5>
    <button type="button" class="btn-close text-reset" on:click={toggleSidebar} aria-label="Close"></button>
  </div>
  <div class="offcanvas-body">
    Menu...
  </div>
</div>
