<script>
  import { createEventDispatcher } from 'svelte';

  export let streak = { current: 0, longest: 0 };
  export let todayReview = null;

  const dispatch = createEventDispatcher();

  function completeReview() {
    dispatch('completeReview');
  }

  $: isCompleted = todayReview && todayReview.status === 'completed';
  $: isSkipped = todayReview && todayReview.status === 'skipped';
</script>

<div class="habit-dashboard">
  <div class="streak-section">
    <div class="streak-display">
      <span class="streak-number">{streak.current}</span>
      <span class="streak-label">day streak</span>
    </div>
    {#if streak.longest > 0}
      <div class="streak-best">
        Best: {streak.longest} days
      </div>
    {/if}
  </div>

  <div class="review-section">
    {#if isCompleted}
      <div class="review-status completed">
        <span class="review-icon">&#10003;</span>
        <span>Today's review done</span>
        {#if todayReview}
          <span class="task-count">
            {todayReview.completedTaskCount}/{todayReview.totalTaskCount} tasks
          </span>
        {/if}
      </div>
    {:else if isSkipped}
      <div class="review-status skipped">
        <span class="review-icon">&#8212;</span>
        <span>Skipped</span>
      </div>
    {:else}
      <button
        class="review-btn"
        on:click={completeReview}
        aria-label="Complete daily review"
      >
        Complete Daily Review
      </button>
    {/if}
  </div>
</div>

<style>
  .habit-dashboard {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .streak-section {
    display: flex;
    align-items: baseline;
    gap: 12px;
  }

  .streak-display {
    display: flex;
    align-items: baseline;
    gap: 4px;
  }

  .streak-number {
    font-size: 24px;
    font-weight: 700;
    color: var(--accent);
  }

  .streak-label {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .streak-best {
    font-size: 11px;
    color: var(--text-secondary);
  }

  .review-section {
    display: flex;
    align-items: center;
  }

  .review-status {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    padding: 4px 10px;
    border-radius: 4px;
  }

  .review-status.completed {
    color: var(--success);
    background: rgba(76, 175, 80, 0.1);
  }

  .review-status.skipped {
    color: var(--text-secondary);
    background: rgba(139, 141, 163, 0.1);
  }

  .review-icon {
    font-size: 14px;
  }

  .task-count {
    color: var(--text-secondary);
    margin-left: 4px;
  }

  .review-btn {
    background: var(--accent);
    color: white;
    border: none;
    padding: 6px 14px;
    border-radius: 4px;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.15s;
  }

  .review-btn:hover {
    background: var(--accent-hover);
  }
</style>
