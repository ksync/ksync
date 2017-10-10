package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCmd_New(t *testing.T) {
	tests := []struct {
		name string
		this *InitCmd
		want *cobra.Command
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &InitCmd{}
		if got := this.New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. InitCmd.New() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestInitCmd_run(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		this *InitCmd
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &InitCmd{}
		this.run(tt.args.cmd, tt.args.args)
	}
}
