<script lang="ts">
  import Login from './lib/Login.svelte';
  import Sidebar from './lib/Sidebar.svelte';
  import { createNote, selectNote } from './lib/stores/notes';
  import { isAuthenticated } from './lib/stores/auth';
  import { toggleSidebar } from './lib/stores/ui';
  import SplitView from './lib/SplitView.svelte';
  import boot from './lib/boot';
  import MediaQuery from './lib/MediaQuery.svelte';

  // TODO
  boot();

  const sidebarId = 'sidebar';
</script>

<main>
  {#if !$isAuthenticated}
    <Login />
  {:else}
    <nav class="navbar navbar-light bg-light p-3 sticky-top">
      <MediaQuery query="(min-width: 680px)" let:matches>
        {#if !matches}
          <button class="btn btn-primary" type="button" on:click={() => selectNote(-1)}>Back</button>
        {/if}
      </MediaQuery>
      <button class="btn btn-primary" type="button" aria-controls={sidebarId} on:click={toggleSidebar}>Menu</button>
      <div class="text-center my-3">
        <button class="btn btn-primary" type="button" on:click={createNote}>New note</button>
      </div>
      <span class="navbar-brand mb-0 h1">Notes</span>
    </nav>

    <Sidebar id={sidebarId} />
    <SplitView />
  {/if}
</main>

<style>
  main {
    display: grid;
    grid-template-areas:
      'nav'
      'content';
    grid-template-rows: auto 1fr;
    height: 100vh;
    width: 100vw;
  }
</style>
