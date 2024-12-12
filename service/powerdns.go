package service

import (
	prefix "NextDomain-Utils/constant"
	"NextDomain-Utils/dto/request"
	"NextDomain-Utils/dto/response"
	"NextDomain-Utils/utils/httpHelper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PowerDNSService struct{}

type IPowerDNSService interface {
	CreateZone(server string, apikey string, req request.ZoneCreate, server_id ...string) (response.ZoneDetail, error)
	DeleteZone(server string, apikey string, req request.ZoneDelete, server_id ...string) error
	GetZones(server string, apikey string, server_id ...string) (response.GetZonesResponse, error)
	GetZonesPdnsAdmin(server string, apikey string, server_id ...string) ([]response.GetZonesPdnsAdminResponse, error)
	GetStatus(server string, apikey string, server_id ...string) (response.GetServerStatusResponse, error)
	ChangeStatus(server string, apikey string, zone_id string, status string, server_id ...string) (response.ZoneChangeStatus, error)
}

func NewIPowerDNSService() IPowerDNSService {
	return &PowerDNSService{}
}

func (p *PowerDNSService) ChangeStatus(server string, apikey string, zone_id string, status string, server_id ...string) (response.ZoneChangeStatus, error) {
	path := fmt.Sprintf("%s/%s/pdnsadmin/zones/%s/%s", server, prefix.Pdns, status, zone_id)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result response.ZoneChangeStatus
	err := httpHelper.ExecuteRequest(ctx, "POST", path, nil, apikey, &result)
	return result, err
}

func (p *PowerDNSService) CreateZone(server string, apikey string, req request.ZoneCreate, server_id ...string) (response.ZoneDetail, error) {
	path := fmt.Sprintf("%s/%s/servers/%s/zones?rrsets=true", server, prefix.Pdns, server_id)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	reqBody, err := json.Marshal(req)
	if err != nil {
		return response.ZoneDetail{}, fmt.Errorf("error marshaling request body: %v", err)
	}
	var result response.ZoneDetail
	err = httpHelper.ExecuteRequest(ctx, "POST", path, bytes.NewReader(reqBody), apikey, &result)
	return result, err
}

func (p *PowerDNSService) DeleteZone(server string, apikey string, req request.ZoneDelete, server_id ...string) error {
	path := fmt.Sprintf("%s/%s/pdnsadmin/zones/%s", server, prefix.Pdns, req.ZoneId)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := httpHelper.ExecuteRequest(ctx, "DELETE", path, nil, apikey, nil)
	return err
}

func (p *PowerDNSService) GetZones(server string, apikey string, server_id ...string) (response.GetZonesResponse, error) {
	path := fmt.Sprintf("%s/%s/servers/%s/zones", server, prefix.Pdns, server_id)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result response.GetZonesResponse
	err := httpHelper.ExecuteRequest(ctx, "GET", path, nil, apikey, &result)
	return result, err
}

func (p *PowerDNSService) GetZonesPdnsAdmin(server string, apikey string, server_id ...string) ([]response.GetZonesPdnsAdminResponse, error) {
	path := fmt.Sprintf("%s/%s/pdnsadmin/zones", server, prefix.Pdns)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result []response.GetZonesPdnsAdminResponse
	err := httpHelper.ExecuteRequest(ctx, "GET", path, nil, apikey, &result)
	return result, err
}

func (p *PowerDNSService) GetStatus(server string, apikey string, server_id ...string) (response.GetServerStatusResponse, error) {
	path := fmt.Sprintf("%s/%s/servers/%s", server, prefix.Pdns, server_id)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result response.GetServerStatusResponse
	err := httpHelper.ExecuteRequest(ctx, "GET", path, nil, apikey, &result)
	return result, err
}
