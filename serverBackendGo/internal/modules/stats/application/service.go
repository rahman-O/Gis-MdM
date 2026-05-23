package application

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/stats/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/stats/port"
)

type Service struct {
	repo port.Repository
}

func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Save(ctx context.Context, stats domain.UsageStats) error {
	if stats.DevicesTotal < 0 {
		stats.DevicesTotal = 0
	}
	if stats.DevicesOnline < 0 {
		stats.DevicesOnline = 0
	}
	return s.repo.Upsert(ctx, stats)
}
