package vault

import (
	"reflect"
	"testing"

	"github.com/ZupIT/go-vault-session/pkg/login"
	"github.com/hashicorp/vault/api"
)

func TestManager_Write(t *testing.T) {
	type fields struct {
		client *api.Client
	}
	type args struct {
		key     string
		data    map[string]interface{}
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{client: buildClient()},
			args:    args{key: "my-test-write", data: buildDummyData()},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{client: buildClient()},
			args:    args{key: "my-test-error", data: buildDummyData(), address: "http://localhost:1234"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.address != "" {
				_ = tt.fields.client.SetAddress(tt.args.address)
			}

			vm := NewVaultManager(tt.fields.client)

			if err := vm.Write(tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Read(t *testing.T) {
	type fields struct {
		client *api.Client
	}
	type args struct {
		key     string
		address string
		data    map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{client: buildClient()},
			args:    args{key: "my-test-write", data: buildDummyData()},
			want:    buildDummyData(),
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  fields{client: buildClient()},
			args:    args{key: "my-test-read-error"},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{client: buildClient()},
			args:    args{key: "my-test-write", address: "http://localhost:1234"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.address != "" {
				_ = tt.fields.client.SetAddress(tt.args.address)
			}

			vm := NewVaultManager(tt.fields.client)

			if tt.args.data != nil {
				vm.Write(tt.args.key, tt.args.data)
			}

			got, err := vm.Read(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	type fields struct {
		client *api.Client
	}
	type args struct {
		data    map[string]interface{}
		key     string
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{client: buildClient()},
			args:    args{data: buildDummyData(), key: "zup"},
			want:    []interface{}{"my-test-list"},
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  fields{client: buildClient()},
			args:    args{key: "notfound"},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{client: buildClient()},
			args:    args{key: "error", address: "http://localhost:1234"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.address != "" {
				_ = tt.fields.client.SetAddress(tt.args.address)
			}

			vm := NewVaultManager(tt.fields.client)

			if tt.args.data != nil {
				_ = vm.Write("zup/my-test-list", buildDummyData())
			}
			got, err := vm.List(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Delete(t *testing.T) {
	type fields struct {
		client *api.Client
	}
	type args struct {
		key     string
		data    map[string]interface{}
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{client: buildClient()},
			args:    args{data: buildDummyData(), key: "test"},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{client: buildClient()},
			args:    args{data: buildDummyData(), key: "test", address: "http://localhost:1234"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.address != "" {
				_ = tt.fields.client.SetAddress(tt.args.address)
			}

			vm := NewVaultManager(tt.fields.client)

			if tt.args.data != nil {
				_ = vm.Write("test", buildDummyData())
			}

			if err := vm.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if data, _ := vm.Read("test"); data != nil {
				t.Errorf("Delete() was not successful, key still there")
			}
		})
	}
}

func buildDummyData() map[string]interface{} {
	dummyData := map[string]interface{}{}
	dummyData["name"] = "git"
	dummyData["password"] = "132465"
	return dummyData
}

func buildClient() *api.Client {
	cfg := api.DefaultConfig()
	_ = cfg.ReadEnvironment()
	c, _ := api.NewClient(cfg)

	l := login.NewHandler(c)
	l.HandleLogin()

	return c
}
