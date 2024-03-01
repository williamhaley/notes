<script lang="ts">
  import Login from './lib/Login.svelte';
  import Menu from './lib/Menu.svelte';
  import { isAuthenticated } from './lib/stores/auth';
  import boot from './lib/boot';
  import MediaQuery from './lib/MediaQuery.svelte';
  import DetailView from './lib/DetailView.svelte';
  import { subscribe, type Note, observeChanges } from './lib/stores/notes';
  import ListView from './lib/ListView.svelte';
  import { onMount } from 'svelte';

  // TODO
  boot();

  const menuId = 'menu';

  let note: Note | undefined;

  subscribe(state => {
    if (state.selectedNoteId) {
      note = state.notes[state.selectedNoteId];
    } else {
      note = undefined;
    }
  });

  onMount(() => {
    observeChanges();
  });
</script>

<main>
  {#if !$isAuthenticated}
    <Login />
  {:else}
    <Menu id={menuId} />
    <MediaQuery query="(min-width: 680px)" let:matches>
      {#if matches}
        <div class="split-view">
          <div class="list">
            <ListView {menuId} />
          </div>
          <div class="detail">
            <DetailView />
          </div>
        </div>
      {:else if note}
        <div class="detail">
          <DetailView />
        </div>
      {:else}
        <div class="list">
          <ListView {menuId} />
        </div>
      {/if}
    </MediaQuery>
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

  .split-view {
    display: flex;
    flex-direction: row;
  }

  .split-view .list {
    width: 30%;
    border-right: 1px solid black;
  }

  .split-view .detail {
    flex: 1;
  }

  .list,
  .detail {
    height: 100vh;
  }
</style>
