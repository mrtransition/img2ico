package sizes

import "fmt"

// Standard sizes for Windows ICO (Vista+ supports up to 256)
var Standard = []int{16, 24, 32, 48, 64, 72, 96, 128, 256}

// Validate checks that all sizes are between 1 and 256.
func Validate(sizes []int) error {
	for _, s := range sizes {
		if s < 1 || s > 256 {
			return fmt.Errorf("size %d is out of range (1-256)", s)
		}
	}
	return nil
}
