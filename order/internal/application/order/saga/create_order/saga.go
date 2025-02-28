package create_order

import "context"

type Saga interface {
	Handle(ctx context.Context, result Result) error
}
