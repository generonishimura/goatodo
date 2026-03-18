<script>
  export let total = 0;
  export let done = 0;
  export let mode = 'normal';
</script>

<div class="status-bar">
  <span class="task-count">{done}/{total} done</span>
  <span class="mode-indicator">{mode.toUpperCase()}</span>
  <span class="shortcuts">
    <kbd>a</kbd> add
    <kbd>j/k</kbd> nav
    <kbd>x</kbd> toggle
    <kbd>e</kbd> edit
    <kbd>d</kbd> delete
    <kbd>1-3</kbd> priority
  </span>
</div>

<style>
  .status-bar {
    display: flex;
    align-items: center;
    padding: 6px 16px;
    background: var(--bg-secondary);
    border-top: 1px solid var(--border);
    font-size: 12px;
    color: var(--text-secondary);
    gap: 16px;
  }

  .task-count {
    color: var(--accent);
  }

  .mode-indicator {
    font-weight: bold;
    color: var(--text-primary);
  }

  .shortcuts {
    margin-left: auto;
    display: flex;
    gap: 8px;
  }

  kbd {
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 3px;
    padding: 1px 4px;
    font-family: monospace;
    font-size: 11px;
  }
</style>
