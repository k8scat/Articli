package gitlab

import "testing"

func TestEncodeProjectID(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "number",
			args: args{
				projectID: "123",
			},
			want: "123",
		},
		{
			name: "namespace with project_path",
			args: args{
				projectID: "namespace/project_path",
			},
			want: "namespace%2Fproject_path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLEncoded(tt.args.projectID); got != tt.want {
				t.Errorf("EncodeProjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}
