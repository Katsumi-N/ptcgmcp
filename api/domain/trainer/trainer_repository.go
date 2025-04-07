package trainer

import (
	"context"
)

type TrainerRepository interface {
	FindById(ctx context.Context, trainerId int) (*Trainer, error)
}
