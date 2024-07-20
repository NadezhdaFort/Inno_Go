package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalSequence(t *testing.T) {
	type args struct {
		mtx [][]int
		ua  []int
	}

	mtx1 := [][]int{
		{0, 2, 3, 0, 0},
		{2, 0, 0, 1, 1},
		{3, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0},
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "mtx 5 verticals 100%",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0, 2},
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "Graph 100%",
			args: args{
				mtx: [][]int{
					{0, 5, 0, 0, 0, 0, 0},
					{5, 0, 15, 50, 35, 0, 0},
					{0, 15, 0, 0, 0, 40, 25},
					{0, 50, 0, 0, 0, 45, 0},
					{0, 35, 0, 0, 0, 0, 30},
					{0, 0, 40, 45, 0, 0, 0},
					{0, 0, 25, 0, 30, 0, 0},
				},
				ua: []int{6, 4, 1, 3, 5, 2},
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "mtx 5 verticals 0%",
			args: args{
				mtx: mtx1,
				ua:  []int{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Empty graph 0%",
			args: args{
				mtx: [][]int{},
				ua:  []int{0},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "mtx 5 verticals 50%",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0},
			},
			want:    50,
			wantErr: false,
		},
		{
			name: "Long user answer 0%",
			args: args{
				mtx: [][]int{
					{0, 2, 0, 0},
					{2, 0, 3, 0},
					{0, 3, 0, 4},
					{0, 0, 4, 0},
				},
				ua: []int{0, 1, 2, 3, 4},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Incorrect graph 0%",
			args: args{
				mtx: [][]int{
					{0, 2, 0, 0},
					{2, 0, 3},
					{0, 3, 0, 4},
					{0, 0, 4, 0, 7},
				},
				ua: []int{0, 1, 2, 3, 4},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid user answer 0%",
			args: args{
				mtx: [][]int{
					{0, 2, 0, 0},
					{2, 0, 3, 0},
					{0, 3, 0, 4},
					{0, 0, 4, 0},
				},
				ua: []int{0, 1, 5, 3},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Non-unique user answer 0%",
			args: args{
				mtx: [][]int{
					{0, 2, 0, 0},
					{2, 0, 3, 0},
					{0, 3, 0, 4},
					{0, 0, 4, 0},
				},
				ua: []int{0, 1, 1, 3},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalSequence(tt.args.mtx, tt.args.ua)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
