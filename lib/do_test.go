package lib_test

import (
	"reflect"
	"testing"

	"github.com/cage1016/alfred-source2target/lib"
	"github.com/cage1016/alfred-source2target/testdata"
)

func TestFdExecute(t *testing.T) {
	type fields struct {
		fn func(cfg lib.DoConfig) []string
	}

	type args struct {
		n map[lib.DoConfig][]string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
	}{
		// {
		// 	name: "ls -ltr",
		// 	prepare: func(f *fields) {
		// 		f.fn = lib.FdExecute
		// 	},
		// 	args: args{
		// 		n: map[lib.DoConfig][]string{
		// 			lib.DoConfig{
		// 				Source:          testdata.Path("target"),
		// 				Arg:             "",
		// 				Type:            "-tf",
		// 				Exclude:         "",
		// 				MaxQueryResults: 100,
		// 			}: []string{
		// 				"./a.txt",
		// 				"./apple-touch-icon.png",
		// 				"./b.txt",
		// 				"./c.txt",
		// 				"./folder1/d.txt",
		// 				"./folder1/e.txt",
		// 				"./folder1/f.txt",
		// 				"./icon-svg.pdf",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "query png",
		// 	prepare: func(f *fields) {
		// 		f.fn = lib.FdExecute
		// 	},
		// 	args: args{
		// 		n: map[lib.DoConfig][]string{
		// 			lib.DoConfig{
		// 				Source:          testdata.Path("target"),
		// 				Arg:             "-e png -e pdf",
		// 				Type:            "-tf",
		// 				Exclude:         "",
		// 				MaxQueryResults: 100,
		// 			}: []string{
		// 				"./apple-touch-icon.png",
		// 				"./icon-svg.pdf",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "query with exclude value",
		// 	prepare: func(f *fields) {
		// 		f.fn = lib.FdExecute
		// 	},
		// 	args: args{
		// 		n: map[lib.DoConfig][]string{
		// 			lib.DoConfig{
		// 				Source:          testdata.Path("target"),
		// 				Arg:             "",
		// 				Type:            "-tf",
		// 				Exclude:         "folder1\n'*.txt'",
		// 				MaxQueryResults: 100,
		// 			}: []string{
		// 				"./apple-touch-icon.png",
		// 				"./icon-svg.pdf",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "query return max result",
		// 	prepare: func(f *fields) {
		// 		f.fn = lib.FdExecute
		// 	},
		// 	args: args{
		// 		n: map[lib.DoConfig][]string{
		// 			lib.DoConfig{
		// 				Source:          testdata.Path("target"),
		// 				Arg:             "",
		// 				Type:            "-tf",
		// 				Exclude:         "",
		// 				MaxQueryResults: 2,
		// 			}: []string{
		// 				"./a.txt",
		// 				"./apple-touch-icon.png",
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "query fail with invalid parameter",
			prepare: func(f *fields) {
				f.fn = lib.FdExecute
			},
			args: args{
				n: map[lib.DoConfig][]string{
					lib.DoConfig{
						Source:          testdata.Path("target"),
						Arg:             "",
						Type:            "-xxx",
						Exclude:         "",
						MaxQueryResults: 2,
					}: []string{
						"error: the argument '--exec <cmd>...' cannot be used with '--exec-batch <cmd>...'",
						"Usage: fd --exec <cmd>... [pattern] [path]...",
						"For more information, try '--help'.",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			for k, v := range tt.args.n {
				if got := f.fn(k); !reflect.DeepEqual(got, v) {
					t.Errorf("FdExecute() = %v, want %v", got, v)
				}
			}
		})
	}
}
