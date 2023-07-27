package scripts

import (
	"github.com/SETTER2000/prove/config"
	"reflect"
	"regexp"
	"testing"
)

func TestCheckEnvironFlag(t *testing.T) {
	type args struct {
		environName string
		flagName    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				environName: "environName",
				flagName:    "flagName",
			},
			want: true,
		}, {
			name: "positive test #2 flag variable set",
			args: args{
				environName: "",
				flagName:    "flagName",
			},
			want: true,
		}, {
			name: "positive test #3 same naming of environment variable and flag",
			args: args{
				environName: "environName",
				flagName:    "DATABASE_URI",
			},
			want: true,
		}, {
			name: "negative test #1",
			args: args{
				environName: "",
				flagName:    "",
			},
			want: false,
		}, {
			name: "negative test #2",
			args: args{
				environName: "environName",
				flagName:    "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckEnvironFlag(tt.args.environName, tt.args.flagName); got != tt.want {
				t.Errorf("CheckEnvironFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHost(t *testing.T) {
	type args struct {
		prove string
		cfg   config.HTTP
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				cfg:   config.HTTP{BaseURL: "http://localhost:8080"},
				prove: "wEuWothteri_t23",
			},
			want: "http://localhost:8080/wEuWothteri_t23",
		}, {
			name: "negative test #1 extra closing slash",
			args: args{
				cfg:   config.HTTP{BaseURL: "http://localhost:8080/"},
				prove: "wEuWothteri_t23",
			},
			want: "http://localhost:8080//wEuWothteri_t23",
		}, {
			name: "negative test #2 missing protocol",
			args: args{
				cfg:   config.HTTP{BaseURL: "localhost:8080"},
				prove: "wEuWothteri_t23",
			},
			want: "localhost:8080/wEuWothteri_t23",
		}, {
			name: "negative test #3 missing protocol and port",
			args: args{
				cfg:   config.HTTP{BaseURL: "localhost"},
				prove: "wEuWothteri_t23",
			},
			want: "localhost/wEuWothteri_t23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHost(tt.args.cfg, tt.args.prove); got != tt.want {
				t.Errorf("GetHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		pattern string
		name    string
		args    args
		want    bool
	}{
		{
			name:    "positive test #1",
			args:    args{n: 10},
			pattern: `[_0-9a-zA-Z]`,
			want:    true,
		}, {
			name:    "positive test #2",
			args:    args{n: 1},
			pattern: `[_0-9a-zA-Z]{3}`,
			want:    true,
		}, {
			name:    "positive test #3",
			args:    args{n: 5},
			pattern: `[_0-9a-zA-Z]{5}`,
			want:    true,
		}, {
			name:    "positive test #4",
			args:    args{n: 0},
			pattern: `[_0-9a-zA-Z]{3}`,
			want:    true,
		}, {
			name:    "negative test #1",
			args:    args{n: -1},
			pattern: `[_0-9a-zA-Z]`,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateString(tt.args.n)
			if match, _ := regexp.MatchString(tt.pattern, got); !match {
				t.Errorf("GenerateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{"\nstring home\n", ""},
			want:    "string home",
			wantErr: false,
		},
		{
			name:    "positive test #2",
			args:    args{"\nstring home:", ":"},
			want:    "string home",
			wantErr: false,
		},
		{
			name:    "positive test #3",
			args:    args{"\nstring home\n", ":"},
			want:    "string home",
			wantErr: false,
		},
		{
			name:    "positive test #4",
			args:    args{"string home", ":"},
			want:    "string home",
			wantErr: false,
		},
		{
			name:    "positive test #5",
			args:    args{"string home", "s"},
			want:    "tring home",
			wantErr: false,
		},
		{
			name:    "positive test #5",
			args:    args{"string home", ""},
			want:    "string home",
			wantErr: false,
		},
		{
			name:    "negative test #1",
			args:    args{"", ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Trim(tt.args.s, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Trim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Trim() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindAllMissingNumbers(t *testing.T) {
	type args struct {
		ar []int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				ar: []int{1, 1, 2, 3, 5, 5, 7, 9, 9, 9},
			},
			want: []int{4, 6, 8},
		}, {
			name: "positive test #2",
			args: args{
				ar: []int{1, 1, 2, 3, 5, 5, 7, 8, 9, 9},
			},
			want: []int{4, 6},
		}, {
			name: "negative test #3",
			args: args{
				ar: []int{1, 1, 5, 6, 7, 8, 9, 9},
			},
			want: []int{2, 3, 4},
		}, {
			name: "negative test #4",
			args: args{
				ar: []int{1, 1, 51, 6, 7, 8, 90, 1, 90, 15, 8, 99, 27, 61, 22, 35, 84, 17, 9},
			},
			want: []int{2, 3, 4, 5, 10, 11, 12, 13, 14, 16, 18, 19, 20, 21, 23, 24, 25, 26, 28, 29, 30, 31, 32, 33, 34, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 52, 53, 54, 55, 56, 57, 58, 59, 60, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 85, 86, 87, 88, 89, 91, 92, 93, 94, 95, 96, 97, 98},
		}, {
			name: "negative test #1",
			args: args{
				ar: []int{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAllMissingNumbers(tt.args.ar)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllMissingNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllMissingNumbers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
