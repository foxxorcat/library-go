package systemutil

func IF[T any](cnod bool, t, f T) T {
	if cnod {
		return t
	}
	return f
}

// IFT 当满足条件时返回t，否则返回空值
func IFT[T any](cnod bool, t T) (f T) {
	if cnod {
		return t
	}
	return f
}

// IFNULL 如果f为空值，则返回t
func IFNULL[T comparable](f T, t T) (cnod T) {
	if f != cnod {
		return f
	}
	return t
}
