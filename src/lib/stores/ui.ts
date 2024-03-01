import { derived, writable } from 'svelte/store';

function createStore() {
  const store = writable({ isMenuOpen: false });

  return {
    toggleMenu() {
      store.update($store => {
        $store.isMenuOpen = !$store.isMenuOpen;
        return $store;
      });
    },
    subscribe: store.subscribe
  };
}

const ui = createStore();

const { subscribe, toggleMenu } = ui;

export { subscribe, toggleMenu };

export const isMenuOpen = derived(ui, store => store.isMenuOpen);
