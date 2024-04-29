package core

import "context"

type BackPressure interface {
	Next(ctx context.Context)
}
