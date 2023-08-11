package str

import "testing"

func TestArToFa(t *testing.T) {
	type args struct {
		ar string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ya", args: args{ar: "ي"}, want: "ی"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArToFa(tt.args.ar); got != tt.want {
				t.Errorf("ArToFa() = %v, want %v", got, tt.want)
			}
		})
	}
}
