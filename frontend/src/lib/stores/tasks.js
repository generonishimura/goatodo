import { writable } from 'svelte/store';

export const tasks = writable([]);
export const selectedIndex = writable(0);
export const isAdding = writable(false);
export const editingId = writable(null);
