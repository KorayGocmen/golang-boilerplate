package slice

import "github.com/google/uuid"

type SliceElem interface {
	string | int | int64 | float64 | uuid.UUID
}

type SliceElemNum interface {
	int | int64 | float64
}

func Contains[SE SliceElem](q SE, list []SE) bool {
	for _, l := range list {
		if q == l {
			return true
		}
	}

	return false
}

func Min[SE SliceElemNum](list ...SE) SE {
	if len(list) == 0 {
		return 0
	}

	min := list[0]
	for _, l := range list {
		if l < min {
			min = l
		}
	}

	return min
}
