package numeric

import "testing"

func TestTransformFa2En(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "persian number",
			input: "۰۹۱۷۱۲۰۹۶۶۸",
			want:  "09171209668",
		},
		{
			name:  "arabic number",
			input: "٠٩١٧١٢٠٩٦٦٨",
			want:  "09171209668",
		},
		{
			name:  "contains persian number",
			input: "یک متن فارسی ۱۲۳ دارد",
			want:  "یک متن فارسی 123 دارد",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformFa2En(tt.input); got != tt.want {
				t.Errorf("TransformFa2En() = %v, want %v", got, tt.want)
			}
		})
	}
}
