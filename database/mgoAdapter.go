package database

import (
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
)

type MgoAdapter struct {
	logger zap.Logger
}

func NewMgoAdapter(url string, logger zap.Logger) (*MgoAdapter, error) {

	return &MgoAdapter{
		logger: logger.With(zap.String("component", "MgoAdapter")),
	}

}
