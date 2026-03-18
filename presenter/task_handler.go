package presenter

import (
	apptask "github.com/i-nishimura/goatodo/application/task"
	domtask "github.com/i-nishimura/goatodo/domain/task"
	"github.com/i-nishimura/goatodo/presenter/dto"
)

type TaskHandler struct {
	createTask   *apptask.CreateTask
	completeTask *apptask.CompleteTask
	listTasks    *apptask.ListTasks
	updateTask   *apptask.UpdateTask
	deleteTask   *apptask.DeleteTask
}

func NewTaskHandler(repo domtask.Repository) *TaskHandler {
	return &TaskHandler{
		createTask:   apptask.NewCreateTask(repo),
		completeTask: apptask.NewCompleteTask(repo),
		listTasks:    apptask.NewListTasks(repo),
		updateTask:   apptask.NewUpdateTask(repo),
		deleteTask:   apptask.NewDeleteTask(repo),
	}
}

type UpdateTaskRequest struct {
	ID       string  `json:"id"`
	Title    *string `json:"title,omitempty"`
	Status   *string `json:"status,omitempty"`
	Priority *int    `json:"priority,omitempty"`
}

type TaskResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *TaskHandler) CreateTask(title string) TaskResponse {
	result := h.createTask.Execute(title)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true, Data: dto.FromTask(result.Value())}
}

func (h *TaskHandler) ListTasks() TaskResponse {
	result := h.listTasks.Execute()
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true, Data: dto.FromTasks(result.Value())}
}

func (h *TaskHandler) CompleteTask(id string) TaskResponse {
	result := h.completeTask.Execute(id)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true}
}

func (h *TaskHandler) UpdateTask(req UpdateTaskRequest) TaskResponse {
	input := apptask.UpdateTaskInput{ID: req.ID}

	if req.Title != nil {
		input.Title = req.Title
	}
	if req.Status != nil {
		s := domtask.Status(*req.Status)
		if !domtask.IsValidStatus(s) {
			return TaskResponse{Success: false, Error: domtask.ErrInvalidStatus}
		}
		input.Status = &s
	}
	if req.Priority != nil {
		p := domtask.Priority(*req.Priority)
		input.Priority = &p
	}

	result := h.updateTask.Execute(input)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true, Data: dto.FromTask(result.Value())}
}

func (h *TaskHandler) DeleteTask(id string) TaskResponse {
	result := h.deleteTask.Execute(id)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true}
}
