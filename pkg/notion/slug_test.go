package notion

import "testing"

func Test_slug(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple slug",
			args: args{
				s: "Hello, world!",
			},
			want: "hello-world",
		},
		{
			name: "with symbols, preserve digits",
			args: args{
				s: "b@r@b-$$$-1223",
			},
			want: "b-r-b-1223",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slug(tt.args.s); got != tt.want {
				t.Errorf("slug() = %v, want %v", got, tt.want)
			}
		})
	}
}
