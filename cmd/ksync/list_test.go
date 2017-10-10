package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestListCmd_New(t *testing.T) {
	type fields struct {
		viper *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   *cobra.Command
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &ListCmd{
			viper: tt.fields.viper,
		}
		if got := this.New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ListCmd.New() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestListCmd_run(t *testing.T) {
	type fields struct {
		viper *viper.Viper
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		this := &ListCmd{
			viper: tt.fields.viper,
		}
		this.run(tt.args.cmd, tt.args.args)
	}
}
