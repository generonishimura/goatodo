<script>
  import { createEventDispatcher } from 'svelte';
  import TaskItem from './TaskItem.svelte';

  export let tasks = [];
  export let selectedIndex = 0;
  export let editingId = null;

  const dispatch = createEventDispatcher();
</script>

<div class="task-list">
  {#if tasks.length === 0}
    <div class="empty-state">
      <p>No tasks yet</p>
      <p class="hint">Press <kbd>a</kbd> to add a task</p>
    </div>
  {:else}
    {#each tasks as task, i (task.id)}
      <TaskItem
        {task}
        selected={i === selectedIndex}
        editing={editingId === task.id}
        on:select={() => dispatch('select', i)}
        on:toggleStatus={() => dispatch('toggleStatus', task)}
        on:update={(e) => dispatch('update', e.detail)}
        on:editDone={() => dispatch('editDone')}
      />
    {/each}
  {/if}
</div>

<style>
  .task-list {
    flex: 1;
    overflow-y: auto;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-secondary);
    gap: 8px;
  }

  .hint {
    font-size: 12px;
  }

  kbd {
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 2px 6px;
    font-size: 12px;
    font-family: monospace;
  }
</style>
