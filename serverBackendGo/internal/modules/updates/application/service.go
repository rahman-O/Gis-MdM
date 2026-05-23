package application

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/updates/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	manifestURL    string
	singleCustomer func(context.Context) (bool, error)
}

func NewService(manifestURL string, singleCustomer func(context.Context) (bool, error)) *Service {
	return &Service{manifestURL: manifestURL, singleCustomer: singleCustomer}
}

func (s *Service) canCheck(ctx context.Context, p *platformauth.Principal) error {
	if p == nil {
		return ErrPermissionDenied
	}
	if p.SuperAdmin {
		return nil
	}
	single, err := s.singleCustomer(ctx)
	if err != nil {
		return err
	}
	if single {
		return nil
	}
	return ErrPermissionDenied
}

func (s *Service) Check(ctx context.Context, p *platformauth.Principal) ([]domain.UpdateEntry, error) {
	if err := s.canCheck(ctx, p); err != nil {
		return nil, err
	}
	body, err := fetchManifest(s.manifestURL)
	if err != nil {
		return nil, err
	}
	var entries []domain.UpdateEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, err
	}
	for i := range entries {
		entries[i].Outdated = entries[i].Version != "" && entries[i].Version != entries[i].CurrentVersion
	}
	return entries, nil
}

func (s *Service) Apply(ctx context.Context, p *platformauth.Principal, req domain.UpdateRequest) ([]domain.UpdateEntry, error) {
	if err := s.canCheck(ctx, p); err != nil {
		return nil, err
	}
	for i := range req.Updates {
		if req.Updates[i].Outdated && !req.Updates[i].UpdateDisabled {
			req.Updates[i].Downloaded = true
			if req.Update {
				req.Updates[i].CurrentVersion = req.Updates[i].Version
				req.Updates[i].Outdated = false
			}
		}
	}
	if req.SendStats {
		// Partial: stats POST stubbed (logged only in Phase 7).
	}
	return req.Updates, nil
}

func fetchManifest(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("manifest fetch failed")
	}
	return io.ReadAll(resp.Body)
}

func SubstituteDomain(manifestURL, host string) string {
	return strings.ReplaceAll(manifestURL, "CUSTOMER_DOMAIN", host)
}
