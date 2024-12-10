package taskmanagementsystem

import (
	"testing"
	"time"
)

func TestTaskManagerSingleton(t *testing.T) {
	manager1 := GetTaskManager()
	manager2 := GetTaskManager()

	if manager1 != manager2 {
		t.Error("TaskManager is not a singleton instance")
	}
}

func TestTaskManager_CreateTask(t *testing.T) {
	manager := GetTaskManager()
	user := NewUser("user1", "Alice", "alice@example.com")
	task := NewTask("task1", "New Task", "Task description", time.Now().Add(24*time.Hour), 1, user)

	manager.CreateTask(task)

	if _, exists := manager.tasks["task1"]; !exists {
		t.Error("Task was not created in TaskManager")
	}
	if len(manager.userTasks[user.GetId()]) != 1 {
		t.Error("Task was not assigned to the user")
	}
}

func TestTaskManager_UpdateTask(t *testing.T) {
	manager := GetTaskManager()
	user1 := NewUser("user1", "Alice", "alice@example.com")
	user2 := NewUser("user2", "Bob", "bob@example.com")
	task := NewTask(
		"task2",
		"Original Task",
		"Original description",
		time.Now().Add(24*time.Hour),
		1,
		user1,
	)

	manager.CreateTask(task)

	updatedTask := NewTask(
		"task2",
		"Updated Task",
		"Updated description",
		time.Now().Add(48*time.Hour),
		2,
		user2,
	)
	manager.UpdateTask(updatedTask)

	// Check updated task details
	updated := manager.tasks["task2"]
	if updated.GetTitle() != "Updated Task" {
		t.Errorf("Expected title to be 'Updated Task', got '%s'", updated.GetTitle())
	}
	if updated.GetAssignedUser().GetId() != "user2" {
		t.Errorf("Expected assigned user to be user2, got '%s'", updated.GetAssignedUser().GetId())
	}
	if len(manager.userTasks[user1.GetId()]) != 0 {
		t.Error("Task was not unassigned from the original user")
	}
	if len(manager.userTasks[user2.GetId()]) != 1 {
		t.Error("Task was not assigned to the new user")
	}
}

func TestTaskManager_DeleteTask(t *testing.T) {
	manager := GetTaskManager()
	user := NewUser("user1", "Alice", "alice@example.com")
	task := NewTask(
		"task3",
		"Task to Delete",
		"Task description",
		time.Now().Add(24*time.Hour),
		1,
		user,
	)

	manager.CreateTask(task)
	manager.DeleteTask(task.GetId())

	if _, exists := manager.tasks[task.GetId()]; exists {
		t.Error("Task was not deleted from TaskManager")
	}
	if len(manager.userTasks[user.GetId()]) != 0 {
		t.Error("Task was not unassigned from the user")
	}
}

func TestTaskManager_SearchTasks(t *testing.T) {
	manager := GetTaskManager()
	user := NewUser("user1", "Alice", "alice@example.com")
	task1 := NewTask(
		"task4",
		"Find me",
		"This is a matching task",
		time.Now().Add(24*time.Hour),
		1,
		user,
	)
	task2 := NewTask("task5", "Unrelated", "Does not match", time.Now().Add(24*time.Hour), 1, user)

	manager.CreateTask(task1)
	manager.CreateTask(task2)

	results := manager.SearchTasks("Find")
	if len(results) != 1 || results[0].GetId() != "task4" {
		t.Error("SearchTasks did not return the expected matching task")
	}
}

func TestTaskManager_FilterTasks(t *testing.T) {
	manager := GetTaskManager()
	now := time.Now()
	user := NewUser("user1", "Alice", "alice@example.com")
	task := NewTask("task6", "Filterable Task", "Task description", now.Add(48*time.Hour), 2, user)
	task.SetStatus(InProgress)

	manager.CreateTask(task)

	results := manager.FilterTasks(InProgress, now, now.Add(72*time.Hour), 2)
	if len(results) != 1 || results[0].GetId() != "task6" {
		t.Error("FilterTasks did not return the expected task")
	}
}

func TestTaskManager_MarkTaskAsCompleted(t *testing.T) {
	manager := GetTaskManager()
	user := NewUser("user1", "Alice", "alice@example.com")
	task := NewTask(
		"task7",
		"Completable Task",
		"Task description",
		time.Now().Add(24*time.Hour),
		1,
		user,
	)

	manager.CreateTask(task)
	manager.MarkTaskAsCompleted(task.GetId())

	if task.GetStatus() != Completed {
		t.Error("MarkTaskAsCompleted did not set the task status to COMPLETED")
	}
}

func TestTaskManager_GetTaskHistory(t *testing.T) {
	manager := GetTaskManager()
	user := NewUser("user1", "Alice", "alice@example.com")
	task1 := NewTask(
		"task8",
		"History Task 1",
		"Task description",
		time.Now().Add(24*time.Hour),
		1,
		user,
	)
	task2 := NewTask(
		"task9",
		"History Task 2",
		"Task description",
		time.Now().Add(48*time.Hour),
		1,
		user,
	)

	manager.CreateTask(task1)
	manager.CreateTask(task2)

	history := manager.GetTaskHistory(user)
	if len(history) != 2 || history[0].GetId() != "task8" || history[1].GetId() != "task9" {
		t.Error("GetTaskHistory did not return the correct task history")
	}
}
