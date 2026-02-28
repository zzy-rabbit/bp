package internal

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/zzy-rabbit/bp/tool/timer/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xtrace"
)

// Register 注册任务（存在则返回错误）
func (s *service) Register(ctx context.Context, name string, spec string, job api.Job) xerror.IError {
	defer xtrace.Trace(ctx)(name)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.taskMap[name]; exists {
		s.ILogger.Error(ctx, "task %s already exists", name)
		return xerror.Extend(xerror.ErrAlreadyExists, "task "+name)
	}

	entryID, err := s.cron.AddFunc(spec, job)
	if err != nil {
		s.ILogger.Error(ctx, "add task %s failed: %v", name, err)
		return xerror.Extend(xerror.ErrInternalError, "add task "+name)
	}

	entry := s.cron.Entry(entryID)

	s.taskMap[name] = &api.Task{
		Name:       name,
		Spec:       spec,
		EntryID:    int(entryID),
		NextRun:    entry.Next,
		CreateTime: time.Now(),
	}
	return nil
}

// Unregister 删除任务
func (s *service) Unregister(ctx context.Context, name string) {
	defer xtrace.Trace(ctx)(name)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, ok := s.taskMap[name]; ok {
		s.cron.Remove(cron.EntryID(task.EntryID))
		delete(s.taskMap, name)
	}
}

// List 返回当前任务列表（只读快照）
func (s *service) List() []api.Task {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]api.Task, 0, len(s.taskMap))
	for _, t := range s.taskMap {
		entry := s.cron.Entry(cron.EntryID(t.EntryID))
		result = append(result, api.Task{
			Name:       t.Name,
			Spec:       t.Spec,
			EntryID:    t.EntryID,
			NextRun:    entry.Next,
			CreateTime: t.CreateTime,
		})
	}
	return result
}
