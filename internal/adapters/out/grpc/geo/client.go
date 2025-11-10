package geo

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var _ ports.GeoClient = &geoClient{}

type geoClient struct {
	conn        *grpc.ClientConn
	pbGeoClient geopb.GeoClient
	timeout     time.Duration
}

func NewClient(host string) (ports.GeoClient, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &geoClient{
		conn:        conn,
		pbGeoClient: geopb.NewGeoClient(conn),
		timeout:     5 * time.Second,
	}, nil
}

func (g geoClient) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	// запрос
	req := &geopb.GetGeolocationRequest{
		Street: street,
	}

	// Отправляем запрос
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	resp, err := g.pbGeoClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	// Создаем и возвращаем Value Object
	return kernel.NewLocation(int(resp.Location.X), int(resp.Location.Y))
}

func (g geoClient) Close() error {
	return g.conn.Close()
}
