package flat

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		log     *slog.Logger
		storage FlatStorage
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(tt.args.log, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
