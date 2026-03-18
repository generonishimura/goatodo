<script>
  import { createEventDispatcher } from 'svelte';

  export let task;
  export let selected = false;
  export let editing = false;

  const dispatch = createEventDispatcher();
  let editTitle = task.title;
  let editInput;

  $: if (editing && editInput) {
    editInput.focus();
  }

  const statusIcons = {
    todo: '○',
    doing: '◐',
    done: '●',
  };

  const priorityLabels = {
    0: '',
    1: 'L',
    2: 'M',
    3: 'H',
  };

  const priorityColors = {
    1: 'var(--priority-low)',
    2: 'var(--priority-medium)',
    3: 'var(--priority-high)',
  };

  function handleEditSubmit() {
    const trimmed = editTitle.trim();
    if (trimmed && trimmed !== task.title) {
      dispatch('update', { id: task.id, title: trimmed });
    }
    dispatch('editDone');
  }

  function handleEditKeydown(e) {
    if (e.key === 'Escape') {
      editTitle = task.title;
      dispatch('editDone');
    }
  }
</script>

<div
  class="task-item"
  class:selected
  class:done={task.status === 'done'}
  on:click={() => dispatch('select')}
>
  <button class="status-btn" on:click|stopPropagation={() => dispatch('toggleStatus')}>
    <span class="status-icon" class:done-icon={task.status === 'done'}>
      {statusIcons[task.status] || '○'}
    </span>
  </button>

  {#if editing}
    <form class="edit-form" on:submit|preventDefault={handleEditSubmit}>
      <input
        bind:this={editInput}
        bind:value={editTitle}
        on:keydown={handleEditKeydown}
        on:blur={handleEditSubmit}
        class="edit-input"
      />
    </form>
  {:else}
    <span class="task-title">{task.title}</span>
  {/if}

  {#if task.priority > 0}
    <span class="priority-badge" style="color: {priorityColors[task.priority]}">
      {priorityLabels[task.priority]}
    </span>
  {/if}
</div>

<style>
  .task-item {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    gap: 10px;
    cursor: default;
    border-left: 3px solid transparent;
    transition: background-color 0.1s;
  }

  .task-item:hover {
    background: var(--bg-hover);
  }

  .task-item.selected {
    background: var(--bg-hover);
    border-left-color: var(--accent);
  }

  .task-item.done .task-title {
    text-decoration: line-through;
    color: var(--text-secondary);
  }

  .status-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 2px;
    font-size: 16px;
    line-height: 1;
    color: var(--text-secondary);
  }

  .status-btn:hover {
    color: var(--accent);
  }

  .done-icon {
    color: var(--success);
  }

  .task-title {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .edit-form {
    flex: 1;
  }

  .edit-input {
    width: 100%;
    padding: 4px 8px;
    background: var(--bg-secondary);
    border: 1px solid var(--accent);
    border-radius: 4px;
    color: var(--text-primary);
    font-size: 14px;
    outline: none;
  }

  .priority-badge {
    font-size: 11px;
    font-weight: bold;
    min-width: 18px;
    text-align: center;
  }
</style>
