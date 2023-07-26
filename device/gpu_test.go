package device

import "testing"

func Test_convertName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "NVIDIA GeForce RTX 2080 Ti",
			args: args{
				name: "NVIDIA GeForce RTX 2080 Ti",
			},
			want: "NVIDIA 2080 Ti",
		},
		{
			name: "GeForce RTX 3090",
			args: args{
				name: "GeForce RTX 3090",
			},
			want: "NVIDIA 3090",
		},
		{
			name: "NVIDIA RTX 2080",
			args: args{
				name: "NVIDIA RTX 2080",
			},
			want: "NVIDIA 2080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertName(tt.args.name); got != tt.want {
				t.Errorf("convertName() = %v, want %v", got, tt.want)
			}
		})
	}
}
