package workerpool

import (
	"context"
)

type TaskID string
type taskType string
type taskMetadata map[string]interface{}

type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)

type TaskDescriptor struct {
	ID       TaskID
	TType    taskType
	Metadata map[string]interface{}
}

type Task struct {
	Descriptor TaskDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
}

type Result struct {
	Value      interface{}
	Err        error
	Descriptor TaskDescriptor
}

func NewTask(descriptor TaskDescriptor, f ExecutionFn, args interface{}) *Task {
	return &Task{
		Descriptor: descriptor,
		ExecFn:     f,
		Args:       args,
	}
}

func (t *Task) process(ctx context.Context) Result {
	value, err := t.ExecFn(ctx, t.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: t.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: t.Descriptor,
	}
}
