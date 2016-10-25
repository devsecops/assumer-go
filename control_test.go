package assumer

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
)

func TestControlPlane_Assume(t *testing.T) {
	type fields struct {
		Plane        Plane
		SerialNumber string
		MfaToken     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *sts.AssumeRoleOutput
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		c := &ControlPlane{
			Plane:        tt.fields.Plane,
			SerialNumber: tt.fields.SerialNumber,
			MfaToken:     tt.fields.MfaToken,
		}
		got, err := c.Assume()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ControlPlane.Assume() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ControlPlane.Assume() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestControlPlane_GetDefaults(t *testing.T) {
	type fields struct {
		Plane        Plane
		SerialNumber string
		MfaToken     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		c := &ControlPlane{
			Plane:        tt.fields.Plane,
			SerialNumber: tt.fields.SerialNumber,
			MfaToken:     tt.fields.MfaToken,
		}
		if err := c.GetDefaults(); (err != nil) != tt.wantErr {
			t.Errorf("%q. ControlPlane.GetDefaults() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCheckMfa(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := CheckMfa(tt.args.token); (err != nil) != tt.wantErr {
			t.Errorf("%q. CheckMfa() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
