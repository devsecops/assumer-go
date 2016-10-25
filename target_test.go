package assumer

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
)

func TestTargetPlane_Assume(t *testing.T) {
	type fields struct {
		Plane Plane
	}
	type args struct {
		c *sts.AssumeRoleOutput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sts.AssumeRoleOutput
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tgt := &TargetPlane{
			Plane: tt.fields.Plane,
		}
		got, err := tgt.Assume(tt.args.c)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. TargetPlane.Assume() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. TargetPlane.Assume() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTargetPlane_GetDefaults(t *testing.T) {
	type fields struct {
		Plane Plane
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tgt := &TargetPlane{
			Plane: tt.fields.Plane,
		}
		if err := tgt.GetDefaults(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TargetPlane.GetDefaults() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
