package todo

import (
	"context"
	"errors"
	"fmt"
	"medium/internal/domain/task"
	"medium/pkg/store"
	"time"
)

func (s *Service) ListTasks(ctx context.Context, status string) (res []task.Response, err error) {
	data, err := s.taskRepository.List(ctx, status)
	if err != nil {
		fmt.Printf("failed to select: %v\n", err)
		return
	}

	res = task.ParseFromEntities(data)

	if isWeekend() {
		for i := range res {
			res[i].Title = "ВЫХОДНОЙ - " + res[i].Title
		}
	}

	return
}
func isWeekend() bool {
	today := time.Now().Weekday()
	return today == time.Saturday || today == time.Sunday
}

func (s *Service) CreateTask(ctx context.Context, req task.Request) (res task.Response, err error) {
	activeAtTime, err := time.Parse("2006-01-02", req.ActiveAt)
	if err != nil {
		return res, fmt.Errorf("failed to parse: %v", err)
	}
	data := task.Entity{
		Title:    &req.Title,
		ActiveAt: &activeAtTime,
	}
	data.ID, err = s.taskRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("failed to create: %v\n", err)
		return
	}
	res = task.ParseFromEntity(data)
	return
}

func (s *Service) GetTask(ctx context.Context, id string) (res task.Response, err error) {
	data, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v\n", err)
	}
	res = task.ParseFromEntity(data)
	return
}

func (s *Service) UpdateTask(ctx context.Context, id string, req task.Request) (err error) {
	activeAtTime, err := time.Parse("2006-01-02", req.ActiveAt)
	if err != nil {
		return fmt.Errorf("failed to parse: %v", err)
	}
	data := task.Entity{
		Title:    &req.Title,
		ActiveAt: &activeAtTime,
	}
	err = s.taskRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, id string) (err error) {
	err = s.taskRepository.Delete(ctx, id)
	if err != nil {
		fmt.Printf("failed to delete by id %v\n", err)
		return
	}
	return
}

func (s *Service) MarkTask(ctx context.Context, id string) (err error) {
	err = s.taskRepository.Mark(ctx, id)
	if err != nil {
		fmt.Printf("failed to mark by id %v\n", err)
		return
	}
	return
}
