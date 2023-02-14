package utils

import "testing"

func TestIsSubfolder(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{"test1", "test2"}, false, false},
		{"test2", args{"test1/test2", "test2"}, false, false},
		{"test3", args{"test1", "test1/test2"}, true, false},
		{"test4", args{"/tmp/test1", "/usr/test2"}, false, false},
		{"test5", args{"/tmp/test1", "/tmp/test1/test2"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSubfolder(tt.args.src, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSubfolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSubfolder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
