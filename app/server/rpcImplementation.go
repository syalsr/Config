package server

import (
	"Config/app/database"
	"Config/proto"
	"context"
	"encoding/json"
	"log"
)

type ConfigWrapper struct {
	proto.UnimplementedConfigWrapperServer
}

func (receiver *ConfigWrapper) GetConfig(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	conn := database.GetConnection().Pool

	getCFG := `
	SELECT 
		service.service, config.version, config.cfg, service.service_id 
	FROM service
	INNER JOIN 
		config ON service.service_id = config.service_id
		AND
		config.version = (SELECT MAX(config.version) FROM config
			WHERE config.service_id = (SELECT service_id FROM service
					WHERE "service" = $1))
        AND
            service.service = $1
	`
	result := &database.Request{}
	err := conn.QueryRow(ctx, getCFG, in.GetService()).Scan(&result.Service, &result.Version, &result.Config, &result.ServiceID)
	if err != nil {
		return nil, err
	}

	log.Printf("Get config for %s", in.GetService())

	return &proto.GetResponse{
			Service: result.Service,
			Version: int32(result.Version),
			Config:  string(result.Config),
		},
		nil
}
func (receiver *ConfigWrapper) CreateConfig(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	conn := database.GetConnection().Pool

	version := int32(1)

	insertService := `
		INSERT INTO service 
		    (service, latest_version) 
		VALUES 
		       ($1, $2) 
		RETURNING service_id
	`

	var serviceID int
	err := conn.QueryRow(ctx, insertService, in.GetConfig().Service, version).Scan(&serviceID)
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
	if err != nil {
		return nil, err
	}
	conn.QueryRow(ctx, insertConfig, serviceID, 1, jsonCFG)

	log.Printf("Create config for %s", in.Config.Service)

	return &proto.CreateResponse{Message: "OK"}, nil
}

func (receiver *ConfigWrapper) DeleteUnusedConfig(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	conn := database.GetConnection().Pool

	CheckForPreviousVersion := `
		SELECT service_id, latest_version FROM service
		WHERE service = $1
	`
	var service_id, latest_version int32
	err := conn.QueryRow(ctx, CheckForPreviousVersion, in.Service).Scan(&service_id, &latest_version)
	if err != nil {
		return nil, err
	}
	log.Println("version", latest_version, in.Version, service_id)
	if latest_version != in.Version {
		DeleteConfig := `
			DELETE FROM config
			WHERE service_id = $1 AND version = $2
		`
		conn.Query(ctx, DeleteConfig, service_id, in.Version)

		log.Printf("Service %s, version %d has been deleted", in.Service, in.Version)

		return &proto.DeleteResponse{Message: "Version has been deleted"}, nil
	}
	log.Printf("Config %s has one version", in.Service)
	return &proto.DeleteResponse{Message: "This version in used"}, nil
}

func (receiver *ConfigWrapper) UpdateConfig(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	conn := database.GetConnection().Pool

	GetServiceIdAndLatestVersion := `
		SELECT service_id, latest_version FROM service
		WHERE service = $1
	`
	var service_id, latest_version int32
	err := conn.QueryRow(ctx, GetServiceIdAndLatestVersion, in.Config.Service).Scan(&service_id, &latest_version)
	if err != nil {
		return nil, err
	}

	newVersion := latest_version + 1

	InsertIntoConfig := `
	INSERT INTO config 
		(service_id, version, cfg) 
	VALUES 
	   ($1, $2, $3) 
	`
	data := in.Config.GetData()
	cfg := database.Config{Key1: data[0].Key1, Key2: data[1].Key2}
	jsonCFG, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	conn.QueryRow(ctx, InsertIntoConfig, service_id, newVersion, jsonCFG)

	UpdateService := `
		UPDATE service
		SET latest_version = $1
		WHERE service_id = $2
	`
	conn.QueryRow(ctx, UpdateService, newVersion, service_id)

	log.Printf("Update %s to %d version", in.Config.Service, newVersion)

	return &proto.CreateResponse{Message: "OK"}, nil
}
