package product

import "fmt"

func parseError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("image service error: %w", err)
}
