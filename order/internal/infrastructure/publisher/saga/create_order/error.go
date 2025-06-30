package create_order

import "fmt"

func parseError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("command not published: %w", err)
}
