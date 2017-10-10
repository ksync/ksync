package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetCmd_New(t *testing.T) {
	tests := []struct {
		name string
		this *GetCmd
		want *cobra.Command
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &GetCmd{}
		if got := this.New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetCmd.New() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetCmd_run(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		this *GetCmd
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &GetCmd{}
		this.run(tt.args.cmd, tt.args.args)
	}
}
