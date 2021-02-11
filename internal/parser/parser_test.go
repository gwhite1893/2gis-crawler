package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				[]byte(`
						<!DOCTYPE html>
						<html>
							<head>
								<title>GOT IT</title>
							</head>
							<body>
								some text
							</body>
						</html>`),
			},
			want: "GOT IT",
		},
		{
			name: "empty tag returns empty string",
			args: args{
				[]byte(`
						<!DOCTYPE html>
						<html>
							<head>
								<title></title>
							</head>
							<body>
								some text
							</body>
						</html>`),
			},
			want: "",
		},
		{
			name: "no title returns nil",
			args: args{
				[]byte(`
						<!DOCTYPE html>
						<html>
							<head></head>
							<body>
								some text
							</body>
						</html>`),
			},
			want: "",
		},
		{
			name: "no title returns nil",
			args: args{
				[]byte(`not valid markdown`),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetTagValue(bytes.NewReader(tt.args.b), "title"), "%v", tt.name)
		})
	}
}
