package services

import (
	"fmt"
	"github.com/SETTER2000/prove/internal/entity"
	"sort"
)

// FindAllMissingNumbers - найти все пропущенные числа.
func FindAllMissingNumbers(data *entity.Data) ([]int, error) {
	var res []int
	ar := data.Ar
	if len(ar) < 1 {
		return nil, fmt.Errorf("error, argument cannot be len: %d", len(ar))
	}
	sort.Slice(ar, func(i, j int) bool {
		return ar[i] < ar[j]
	})

	registry := make([]int, ar[len(ar)-1]+1)
	for i, v := range ar {
		registry[v] = i
	}
	for i := 0; i < len(registry); i++ {
		if registry[i] == 0 {
			res = append(res, i)
		}
	}
	return res[1:], nil
}
