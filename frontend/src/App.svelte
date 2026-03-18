<script>
  import { onMount, tick } from 'svelte';
  import TaskList from './lib/components/TaskList.svelte';
  import AddTask from './lib/components/AddTask.svelte';
  import StatusBar from './lib/components/StatusBar.svelte';

  let tasks = [];
  let selectedIndex = 0;
  let isAdding = false;
  let editingId = null;
  let addTaskComponent;
  let mode = 'normal';

  // These will be provided by Wails bindings
  let api = null;

  onMount(async () => {
    // Dynamic import of Wails bindings (generated at build time)
    try {
      const mod = await import('../wailsjs/go/presenter/TaskHandler.js');
      api = mod;
      await loadTasks();
    } catch (e) {
      console.warn('Wails bindings not available (dev mode?):', e);
    }
  });

  async function loadTasks() {
    if (!api) return;
    const res = await api.ListTasks();
    if (res.success) {
      tasks = res.data || [];
    }
  }

  async function addTask(e) {
    if (!api) return;
    const { title } = e.detail;
    const res = await api.CreateTask(title);
    if (res.success) {
      await loadTasks();
      isAdding = false;
      mode = 'normal';
    }
  }

  async function toggleStatus(e) {
    if (!api) return;
    const task = e.detail;
    let nextStatus;
    if (task.status === 'todo') nextStatus = 'doing';
    else if (task.status === 'doing') nextStatus = 'done';
    else return;

    const res = await api.UpdateTask({ id: task.id, status: nextStatus });
    if (res.success) {
      await loadTasks();
    }
  }

  async function updateTask(e) {
    if (!api) return;
    const { id, title } = e.detail;
    const res = await api.UpdateTask({ id, title });
    if (res.success) {
      await loadTasks();
    }
  }

  async function deleteSelectedTask() {
    if (!api || tasks.length === 0) return;
    const task = tasks[selectedIndex];
    if (!task) return;
    const res = await api.DeleteTask(task.id);
    if (res.success) {
      await loadTasks();
      if (selectedIndex >= tasks.length) {
        selectedIndex = Math.max(0, tasks.length - 1);
      }
    }
  }

  async function setPriority(priority) {
    if (!api || tasks.length === 0) return;
    const task = tasks[selectedIndex];
    if (!task) return;
    const res = await api.UpdateTask({ id: task.id, priority });
    if (res.success) {
      await loadTasks();
    }
  }

  function handleKeydown(e) {
    // Don't handle keys when adding or editing
    if (isAdding || editingId) return;

    switch (e.key) {
      case 'j':
        e.preventDefault();
        if (selectedIndex < tasks.length - 1) selectedIndex++;
        break;
      case 'k':
        e.preventDefault();
        if (selectedIndex > 0) selectedIndex--;
        break;
      case 'a':
        e.preventDefault();
        isAdding = true;
        mode = 'insert';
        tick().then(() => {
          if (addTaskComponent) addTaskComponent.focus();
        });
        break;
      case 'x':
        e.preventDefault();
        if (tasks.length > 0) {
          toggleStatus({ detail: tasks[selectedIndex] });
        }
        break;
      case 'e':
      case 'Enter':
        e.preventDefault();
        if (tasks.length > 0) {
          editingId = tasks[selectedIndex].id;
          mode = 'edit';
        }
        break;
      case 'd':
        e.preventDefault();
        deleteSelectedTask();
        break;
      case '1':
        e.preventDefault();
        setPriority(1);
        break;
      case '2':
        e.preventDefault();
        setPriority(2);
        break;
      case '3':
        e.preventDefault();
        setPriority(3);
        break;
      case '0':
        e.preventDefault();
        setPriority(0);
        break;
      case 'G':
        e.preventDefault();
        if (tasks.length > 0) selectedIndex = tasks.length - 1;
        break;
      case 'g':
        e.preventDefault();
        if (tasks.length > 0) selectedIndex = 0;
        break;
    }
  }

  function handleCancelAdd() {
    isAdding = false;
    mode = 'normal';
  }

  function handleEditDone() {
    editingId = null;
    mode = 'normal';
  }

  function handleSelect(e) {
    selectedIndex = e.detail;
  }

  $: doneCount = tasks.filter(t => t.status === 'done').length;
</script>

<svelte:window on:keydown={handleKeydown} />

<main>
  <header>
    <h1>goatodo</h1>
  </header>

  {#if isAdding}
    <AddTask
      bind:this={addTaskComponent}
      on:add={addTask}
      on:cancel={handleCancelAdd}
    />
  {/if}

  <TaskList
    {tasks}
    {selectedIndex}
    {editingId}
    on:select={handleSelect}
    on:toggleStatus={toggleStatus}
    on:update={updateTask}
    on:editDone={handleEditDone}
  />

  <StatusBar total={tasks.length} done={doneCount} {mode} />
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  header {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    -webkit-app-region: drag;
  }

  h1 {
    font-size: 16px;
    font-weight: 600;
    color: var(--accent);
    letter-spacing: 0.5px;
  }
</style>
