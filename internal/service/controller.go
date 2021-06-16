package service

import "net/url"

type Service struct {
	proxyPool            []*url.URL
	rusprofileController *RusprofileService
}

func New(pool []*url.URL) *Service {
	return &Service{
		proxyPool: pool,
	}
}

func (s *Service) RusprofileController() *RusprofileService {
	if s.rusprofileController != nil {
		return s.rusprofileController
	}

	s.rusprofileController = &RusprofileService{service: s}

	return s.rusprofileController
}
