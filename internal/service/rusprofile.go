package service

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/html"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"test_task/pkg/api"
	"test_task/pkg/models"
	"time"
)

type RusprofileService struct {
	api.UnimplementedRusprofileServiceServer
	service *Service
}

func (rs *RusprofileService) GetData(ctx context.Context, req *api.RequestPersonalData) (*api.ResponsePersonalData, error) {
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	if len(req.Inn) != 10 || strings.IndexFunc(req.Inn, isNotDigit) != -1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid inn, inn is a 10-digit personal number")
	}

	client := &http.Client{Timeout: time.Second * 5}
	var resp *http.Response
	var err, errS error
	proxyPollLen := len(rs.service.proxyPool)
	for i := 0; i < proxyPollLen+1; i++ {
		if proxyPollLen > 0 && i > 0 {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(rs.service.proxyPool[rand.Intn(proxyPollLen)])}
		}
		u := url.URL{
			Scheme:   "https",
			Host:     "www.rusprofile.ru",
			Path:     "search",
			RawQuery: fmt.Sprintf("query=%s", req.Inn),
		}
		resp, err = client.Get(u.String())
		if err != nil {
			errS = status.Errorf(codes.Aborted, err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			if errS != nil {
				errS = nil
			}
			break
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			if errS != nil {
				errS = nil
			}
			resp.Body.Close()
			continue
		} else {
			resp.Body.Close()
			break
		}
	}
	if errS != nil {
		return nil, status.Errorf(codes.Aborted, errS.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, status.Errorf(codes.Aborted, "service doesn't work")
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		panic(err)
	}
	if strings.Contains(buf.String(), "было найдено 0 результатов на портале Rusprofile.ru") {
		return nil, status.Errorf(codes.NotFound, "company not found")
	}

	node, err := html.Parse(buf)
	if err != nil {
		panic(err)
	}

	Rusprofile, err := models.NewRusprofileFromWEB(node)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return &api.ResponsePersonalData{
		Inn: Rusprofile.INN,
		Kpp: Rusprofile.KPP,
		Ceo: Rusprofile.CEO,
	}, nil
}
