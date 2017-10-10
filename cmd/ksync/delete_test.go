package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestDeleteCmd_New(t *testing.T) {
	tests := []struct {
		name string
		this *DeleteCmd
		want *cobra.Command
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &DeleteCmd{}
		if got := this.New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DeleteCmd.New() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDeleteCmd_run(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		this *DeleteCmd
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &DeleteCmd{}
		this.run(tt.args.cmd, tt.args.args)
	}
}
