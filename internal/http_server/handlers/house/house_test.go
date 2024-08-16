package house

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func TestHappyCreate(t *testing.T) {
}

func TestCreate(t *testing.T) {
	type args struct {
		log     *slog.Logger
		storage HouseStorage
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

func TestFlats(t *testing.T) {
	type args struct {
		log     *slog.Logger
		storage HouseStorage
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
			if got := Flats(tt.args.log, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flats() = %v, want %v", got, tt.want)
			}
		})
	}
}
