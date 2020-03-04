package main

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestRad(t *testing.T) {
	type args struct {
		d float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rad(tt.args.d); got != tt.want {
				t.Errorf("Rad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeg(t *testing.T) {
	type args struct {
		r float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Deg(tt.args.r); got != tt.want {
				t.Errorf("Deg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRouterCoordinates(t *testing.T) {
	type args struct {
		debug       bool
		routerArray [1000]Router
		routerName  string
	}
	tests := []struct {
		name  string
		args  args
		want  float32
		want1 float32
		want2 float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getRouterCoordinates(tt.args.debug, tt.args.routerArray, tt.args.routerName)
			if got != tt.want {
				t.Errorf("getRouterCoordinates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getRouterCoordinates() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getRouterCoordinates() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
