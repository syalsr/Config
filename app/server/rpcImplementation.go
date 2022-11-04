package server

import (
	"Config/app/database"
	"Config/app/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"
)


type ConfigWrapper struct {
	proto.UnimplementedConfigWrapperServer
}

func (receiver *ConfigWrapper) Get(ctx context.Context, in *proto.Service) (*proto.Data, error) {
	fmt.Println("Get")
	conn := database.GetConnection().Pool

	getCFG := `
	SELECT 
		service.service, config.version, config.cfg, service.service_id 
	FROM service
	INNER JOIN 
		config ON service.service_id = config.service_id
		AND
		config.version = (SELECT config.version FROM config
			WHERE config.service_id = (SELECT service_id FROM service
					WHERE "service" = $1))
        AND
            service.service = $1
	`
	result := &database.Request{}
	fmt.Println(in.GetService())
	err := conn.QueryRow(ctx, getCFG, in.GetService()).Scan(&result.Service, &result.Version, &result.Config)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	log.Printf("Get config for %s", in.GetService())

	return nil, nil
}
func (receiver *ConfigWrapper) Create(ctx context.Context, in *proto.Config) (*proto.Service, error) {
	conn := database.GetConnection().Pool

	insertService := `
		INSERT INTO service 
		    (service, latest_version) 
		VALUES 
		       ($1, $2) 
		RETURNING service_id
	`

	var serviceID int
	err := conn.QueryRow(ctx, insertService, in.GetConfig().Service, 1).Scan(&serviceID)
	if err != nil {
		return nil, err
	}

	insertConfig := `
		INSERT INTO config 
		    (service_id, version, cfg) 
		VALUES 
		       ($1, $2, $3) 
	`
	data := in.Config.GetData()
	cfg := database.Config{Key1: data[0].Key1, Key2: data[1].Key2}
	jsonCFG, err := json.Marshal(cfg)
	if err != nil{
		return nil, err
	}
	conn.QueryRow(ctx, insertConfig, serviceID, 1, jsonCFG)
	
	log.Printf("Create config for %s", in.Config.Service)

	return &proto.Service{Service: in.Config.Service}, nil
}

func (receiver *ConfigWrapper) Delete(context.Context, *proto.Service) (*proto.Service, error) {
	return nil, nil
}

func (receiver *ConfigWrapper) Update(context.Context, *proto.Service) (*proto.Service, error) {
	

	return nil, nil
}