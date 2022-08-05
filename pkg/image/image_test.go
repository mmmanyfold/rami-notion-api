package image

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
			gotW, gotH, err := Size(tt.args.imgPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW != tt.wantW {
				t.Errorf("Size() gotW = %v, want %v", gotW, tt.wantW)
			}
			if gotH != tt.wantH {
				t.Errorf("Size() gotH = %v, want %v", gotH, tt.wantH)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		args     args
		wantFile *os.File
		wantErr  bool
	}{
		{
			name: "download image from internet",
			args: args{
				url: "https://i.imgur.com/m1UIjW1.jpg",
			},
			wantFile: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Download(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotFile, tt.wantFile) {
			//	t.Errorf("Download() gotFile = %v, want %v", gotFile, tt.wantFile)
			//}
		})
	}
}
