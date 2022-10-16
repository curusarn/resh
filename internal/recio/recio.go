package recio

import (
	"go.uber.org/zap"
)

type RecIO struct {
	sugar *zap.SugaredLogger
}

func New(sugar *zap.SugaredLogger) RecIO {
	return RecIO{sugar: sugar}
}
