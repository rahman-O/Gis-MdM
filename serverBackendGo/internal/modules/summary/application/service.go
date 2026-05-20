package application

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/summary/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary/port"
)

// Service provides summary use cases.
type Service struct {
	repo port.SummaryRepository
}

func NewService(repo port.SummaryRepository) *Service {
	return &Service{repo: repo}
}

// GetDeviceStats returns dashboard statistics for the authenticated tenant.
func (s *Service) GetDeviceStats(ctx context.Context, customerID int, userID int64) (*domain.DeviceStats, error) {
	stats, err := s.repo.GetDeviceStats(ctx, customerID, userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return domain.EmptyDeviceStats(), nil
	}
	if stats.StatusSummary == nil || stats.InstallSummary == nil {
		return domain.EmptyDeviceStats(), nil
	}
	return stats, nil
}
