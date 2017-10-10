package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		main()
	}
}

func Test_initPersistent(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		initPersistent(tt.args.cmd, tt.args.args)
	}
}

func Test_initClient(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		initClient()
	}
}
