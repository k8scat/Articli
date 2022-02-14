package gitlab

import (
	"testing"
)

func TestGetFile(t *testing.T) {
	type args struct {
		projectId string
		filePath  string
		ref       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get file",
			args: args{
				projectId: "15242",
				filePath:  "README.md",
				ref:       "main",
			},
			want:    "README.md",
			wantErr: false,
		},
		{
			name: "Get file with namespace/project_name",
			args: args{
				projectId: "kube/kube",
				filePath:  "README.md",
				ref:       "main",
			},
			want:    "README.md",
			wantErr: false,
		},
		{
			name: "Get file without ref",
			args: args{
				projectId: "kube/kube",
				filePath:  "README.md",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Get deep file",
			args: args{
				projectId: "kube/kube",
				filePath:  "test12/test",
				ref:       "main",
			},
			want:    "test",
			wantErr: false,
		},
	}
	setupClient()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetFile(tt.args.projectId, tt.args.filePath, tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" && got.FileName != tt.want {
				t.Errorf("GetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
