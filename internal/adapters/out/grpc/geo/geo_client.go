package geo

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var _ ports.GeoLocationGateway = &geoLocationService{}

type geoLocationService struct {
	conn    *grpc.ClientConn
	client  geopb.GeoClient
	timeout time.Duration
}

func NewGeoLocationService(host string) (*geoLocationService, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbClient := geopb.NewGeoClient(conn)

	return &geoLocationService{
		conn:    conn,
		client:  pbClient,
		timeout: 5 * time.Second,
	}, nil
}

func (g *geoLocationService) Close() error {
	return g.conn.Close()
}

func (g *geoLocationService) DefineLocation(ctx context.Context, street string) (kernel.Location, error) {
	// Формируем запрос
	req := &geopb.GetGeolocationRequest{
		Street: street,
	}

	// Делаем запрос
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	resp, err := g.client.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	// Создаем и возвращаем VO Geo
	location, err := kernel.NewLocation(uint8(resp.Location.X), uint8(resp.Location.Y))
	if err != nil {
		return kernel.Location{}, err
	}
	return location, nil
}
