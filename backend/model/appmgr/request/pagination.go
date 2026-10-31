package request

func NormalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
