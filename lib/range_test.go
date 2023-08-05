package lib

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageRangeRegex(t *testing.T) {

	type fields struct {
		c *regexp.Regexp
	}

	type args struct {
		n map[string]string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
	}{
		{
			name: "test",
			prepare: func(f *fields) {
				f.c = PageRangeRegex
			},
			args: args{
				n: map[string]string{
					"#3-4,5-9,-5":                      "#3-4,5-9,-5",
					"#3-4,5-9,4-,333d,-":               "#3-4,5-9,4-,333",
					"#33-45,333":                       "#33-45,333",
					"#35":                              "#35",
					"#-2":                              "#-2",
					"#1-4":                             "#1-4",
					"#5-":                              "#5-",
					"#5-10":                            "#5-10",
					"#1,1-2,-1,4-,-":                   "#1,1-2,-1,4-,-",
					"#1,-2,3-6,5-,-,dfhfjfdj":          "#1,-2,3-6,5-,-,",
					"#3":                               "#3",
					"#-":                               "#-",
					"--changed-within 1week #3-4,5":    "#3-4,5",
					"--changed-within 1week #3-4,5 dd": "#3-4,5",
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
				if got := f.c.FindAllString(k, -1); got[0] != v {
					t.Errorf("FindAllString = %v, want %v", got, v)
				}
			}
		})
	}
}

func TestParseRangeNumber(t *testing.T) {

	type res struct {
		res []Range
		err error
	}

	type fields struct {
		c func(string, int) (*[]Range, error)
	}

	type args struct {
		maxPage int
		n       map[string]res
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
	}{
		{
			name: "test success",
			prepare: func(f *fields) {
				f.c = ParseRangeNumber
			},
			args: args{
				n: map[string]res{
					"3-4,5-9,-5,5-": {
						res: []Range{
							{Start: 3, End: 4},
							{Start: 5, End: 9},
							{Start: 1, End: 5},
							{Start: 5, End: 10},
						},
						err: nil,
					},
				},
				maxPage: 10,
			},
		},
		{
			name: "test fail with invalid page number",
			prepare: func(f *fields) {
				f.c = ParseRangeNumber
			},
			args: args{
				n: map[string]res{
					"3-4,20": {
						res: []Range{},
						err: fmt.Errorf("page range \"20\" is out of total page \"10\""),
					},
				},
				maxPage: 10,
			},
		},
		{
			name: "test fail with start > end",
			prepare: func(f *fields) {
				f.c = ParseRangeNumber
			},
			args: args{
				n: map[string]res{
					"8-4": {
						res: []Range{},
						err: fmt.Errorf("page range start \"8\" must be less or equal to end \"4\""),
					},
				},
				maxPage: 10,
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
				got, err := f.c(k, tt.args.maxPage)
				if err != nil {
					assert.Equal(t, v.err, err)
				} else {
					assert.Equal(t, v.res, *got)
				}
			}
		})
	}
}

func TestIsPageRangeValid(t *testing.T) {
	type res struct {
		res bool
	}

	type fields struct {
		c func(string) bool
	}

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#3,4,5 is valid",
			args: args{
				s: "#3,4,5",
			},
			want: true,
		},
		{
			name: "#3-4,5-9,-5 is valid",
			args: args{
				s: "#3-4,5-9,-5",
			},
			want: true,
		},
		{
			name: "-d, is invalid",
			args: args{
				s: "-d,",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPageRangeValid(tt.args.s); got != tt.want {
				t.Errorf("IsPageRangeValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRanges_IsInRange(t *testing.T) {
	type res struct {
		res bool
	}

	type fields struct {
		fn func(int) bool
	}

	type args struct {
		val int
	}
	tests := []struct {
		name string
		r    Ranges
		args args
		want bool
	}{
		{
			name: "2 in range",
			r: Ranges{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
			args: args{
				val: 2,
			},
			want: true,
		},
		{
			name: "4 not in range",
			r: Ranges{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
			args: args{
				val: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsInRange(tt.args.val); got != tt.want {
				t.Errorf("Ranges.IsInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
