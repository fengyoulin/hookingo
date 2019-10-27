package hookingo

import "errors"

var (
	// ErrRelativeAddr just as its description
	ErrRelativeAddr = errors.New("relative address in instruction")
)

func ensureLength(src []byte, size int) (int, error) {
	var l int
	for l < size {
		il, err := analysis(src)
		if err != nil {
			return 0, err
		}
		l += il
		src = src[il:]
	}
	return l, nil
}
