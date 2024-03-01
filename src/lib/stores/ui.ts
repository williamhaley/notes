import { derived, writable } from 'svelte/store';

function createStore() {
  const store = writable({ isSidebarOpen: false });

  return {
    toggleSidebar() {
      store.update($store => {
        $store.isSidebarOpen = !$store.isSidebarOpen;
        return $store;
      });
    },
    subscribe: store.subscribe
  };
}

const ui = createStore();

const { subscribe, toggleSidebar } = ui;

export { subscribe, toggleSidebar };

export const isSidebarOpen = derived(ui, store => store.isSidebarOpen);
