package random_code

import "testing"

func TestGenerateRandomDigit(t *testing.T) {
	type args struct {
		length uint
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "4 digit", args: args{length: 4}},
		{name: "6 digit", args: args{length: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRandomDigit(tt.args.length)
			if uint(len(got)) != tt.args.length {
				t.Errorf("GenerateRandomDigit() = %v, want %v", len(got), tt.args.length)
			}
		})
	}
}
