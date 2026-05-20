package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	repo *postgres.AuditRepository
}

func NewService(repo *postgres.AuditRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, f domain.AuditLogFilter) (domain.PaginatedAudit, error) {
	if p == nil || !p.CanPluginAuditAccess() {
		return domain.PaginatedAudit{}, ErrPermissionDenied
	}
	items, total, err := s.repo.Search(ctx, int64(p.CustomerID), f)
	if err != nil {
		return domain.PaginatedAudit{}, err
	}
	return domain.PaginatedAudit{Items: items, Total: total}, nil
}
