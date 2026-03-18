<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();
  let title = '';
  let inputEl;

  export function focus() {
    if (inputEl) inputEl.focus();
  }

  function handleSubmit() {
    const trimmed = title.trim();
    if (trimmed) {
      dispatch('add', { title: trimmed });
      title = '';
    }
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') {
      title = '';
      dispatch('cancel');
    }
  }
</script>

<form class="add-task" on:submit|preventDefault={handleSubmit}>
  <input
    bind:this={inputEl}
    bind:value={title}
    on:keydown={handleKeydown}
    placeholder="New task... (Enter to add, Esc to cancel)"
    aria-label="New task title"
    class="add-input"
    autofocus
  />
</form>

<style>
  .add-task {
    padding: 8px 16px;
    border-bottom: 1px solid var(--border);
  }

  .add-input {
    width: 100%;
    padding: 10px 12px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 6px;
    color: var(--text-primary);
    font-size: 14px;
    outline: none;
  }

  .add-input:focus {
    border-color: var(--accent);
  }

  .add-input::placeholder {
    color: var(--text-secondary);
  }
</style>
