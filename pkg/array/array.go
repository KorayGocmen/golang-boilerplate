package array

func BeginEndInclude(start, stop int64) []int64 {
	var list []int64
	for i := start; i <= stop; i++ {
		list = append(list, i)
	}
	return list
}
