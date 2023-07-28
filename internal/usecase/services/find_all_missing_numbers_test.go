package services

import (
	"github.com/SETTER2000/prove/internal/entity"
	"reflect"
	"testing"
)

type Data struct {
	Ar []int
}

func TestFindAllMissingNumbers(t *testing.T) {
	type args struct {
		data *entity.Data
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
				data: &entity.Data{Ar: []int{1, 1, 2, 3, 5, 5, 7, 9, 9, 9}},
			},
			want: []int{4, 6, 8},
		}, {
			name: "positive test #2",
			args: args{
				data: &entity.Data{Ar: []int{1, 1, 2, 3, 5, 5, 7, 8, 9, 9}},
			},
			want: []int{4, 6},
		}, {
			name: "negative test #3",
			args: args{
				data: &entity.Data{Ar: []int{1, 1, 5, 6, 7, 8, 9, 9}},
			},
			want: []int{2, 3, 4},
		}, {
			name: "negative test #4",
			args: args{
				data: &entity.Data{Ar: []int{1, 1, 51, 6, 7, 8, 90, 1, 90, 15, 8, 99, 27, 61, 22, 35, 84, 17, 9}},
			},
			want: []int{2, 3, 4, 5, 10, 11, 12, 13, 14, 16, 18, 19, 20, 21, 23, 24, 25, 26, 28, 29, 30, 31, 32, 33, 34, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 52, 53, 54, 55, 56, 57, 58, 59, 60, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 85, 86, 87, 88, 89, 91, 92, 93, 94, 95, 96, 97, 98},
		}, {
			name: "negative test #1",
			args: args{
				data: &entity.Data{Ar: []int{}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAllMissingNumbers(tt.args.data)
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
