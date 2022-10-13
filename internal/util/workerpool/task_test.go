package workerpool

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	descriptor = TaskDescriptor{
		ID:    TaskID("1"),
		TType: taskType("anyType"),
		Metadata: taskMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}
	execFn = func(ctx context.Context, args interface{}) (interface{}, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}
)

func Test_task_Process(t *testing.T) {
	ctx := context.TODO()

	type fields struct {
		descriptor TaskDescriptor
		execFn     ExecutionFn
		args       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Result
	}{
		{
			name: "task execution success",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       10,
			},
			want: Result{
				Value:      20,
				Descriptor: descriptor,
			},
		},
		{
			name: "task execution failure",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       "10",
			},
			want: Result{
				Err:        errDefault,
				Descriptor: descriptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				Descriptor: tt.fields.descriptor,
				ExecFn:     tt.fields.execFn,
				Args:       tt.fields.args,
			}

			got := task.process(ctx)
			if tt.want.Err != nil {
				if !cmp.Equal(got.Err, tt.want.Err, cmpopts.EquateErrors()) {
					t.Errorf("execute() = %v, wantError %v", got.Err, tt.want.Err)
				}
				return
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
