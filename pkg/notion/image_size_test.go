package notion

import (
	"os"
	"path"
	"testing"
)

func Test_imageSize(t *testing.T) {
	cwd, _ := os.Getwd()
	testImgPath := path.Join(cwd, "../..", "./fixtures", "rami.jpg")

	type args struct {
		imgPath string
	}
	tests := []struct {
		name    string
		args    args
		wantW   int
		wantH   int
		wantErr bool
	}{{
		name: "image size jpg",
		args: args{
			imgPath: testImgPath,
		},
		wantW:   1294,
		wantH:   2000,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH, err := imageSize(tt.args.imgPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW != tt.wantW {
				t.Errorf("imageSize() gotW = %v, want %v", gotW, tt.wantW)
			}
			if gotH != tt.wantH {
				t.Errorf("imageSize() gotH = %v, want %v", gotH, tt.wantH)
			}
		})
	}
}
