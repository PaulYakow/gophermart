package workerpool

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	tasksCount  = 20
	workerCount = 3
)

func TestWorkerPool(t *testing.T) {
	wp := NewPool(workerCount, workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.GenerateFrom(testTasks()...)

	go wp.RunBackground(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				return
			}

			i, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				t.Fatalf("unexpected error with id = %v: %v", r.Descriptor.ID, err)
			}

			val := r.Value.(int)
			if val != int(i)*2 {
				t.Fatalf("wrong value %v; expected %v", val, int(i)*2)
			}

			if i == tasksCount-1 {
				wp.Stop()
			}

		case <-ctx.Done():
			return
		default:
		}
	}
}

func TestWorkerPool_TimeOut(t *testing.T) {
	wp := NewPool(workerCount, workerCount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond*10)
	defer cancel()

	go wp.RunBackground(ctx)

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.DeadlineExceeded {
				t.Fatalf("expected error: %v; got: %v", context.DeadlineExceeded, r.Err)
			}
		case <-ctx.Done():
			return
		default:
		}
	}
}

func TestWorkerPool_Cancel(t *testing.T) {
	wp := NewPool(workerCount, workerCount)

	ctx, cancel := context.WithCancel(context.TODO())

	go wp.RunBackground(ctx)
	cancel()

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.Canceled {
				t.Fatalf("expected error: %v; got: %v", context.Canceled, r.Err)
			}
		case <-ctx.Done():
			return
		default:
		}
	}
}

func testTasks() []Task {
	tasks := make([]Task, tasksCount)
	for i := 0; i < tasksCount; i++ {
		tasks[i] = Task{
			Descriptor: TaskDescriptor{
				ID:       TaskID(fmt.Sprintf("%v", i)),
				TType:    "anyType",
				Metadata: nil,
			},
			ExecFn: execFn,
			Args:   i,
		}
	}
	return tasks
}
