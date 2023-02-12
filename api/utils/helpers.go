package utils

func Repeat(times int, do func(at int)) {
	for i := 0; i < times; i += 1 {
		do(i)
	}
}

func ForEach[T interface{}](slice *[]T, do func(with T)) {
	for _, with := range *slice {
		do(with)
	}
}