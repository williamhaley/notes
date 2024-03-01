<script lang="ts">
  import { Offcanvas } from 'bootstrap';
  import { onMount } from 'svelte';
  import { subscribe, toggleMenu } from './stores/ui';

  export let id: string;

  let menu: Offcanvas;
  let menuElement: HTMLDivElement;

  onMount(() => {
    menu = new Offcanvas(menuElement);
  });

  subscribe(({ isMenuOpen }) => {
    if (!menu) {
      return;
    }

    if (isMenuOpen) {
      menu.show();
    } else {
      menu.hide();
    }
  });
</script>

<div bind:this={menuElement} class="offcanvas offcanvas-start" tabindex="-1" {id} aria-labelledby="menuLabel">
  <div class="offcanvas-header">
    <h5 id="menuLabel">Menu</h5>
    <button type="button" class="btn-close text-reset" on:click={toggleMenu} aria-label="Close"></button>
  </div>
  <div class="offcanvas-body">Menu...</div>
</div>

<style>
  .offcanvas {
    padding: env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom) env(safe-area-inset-left);
  }
</style>
