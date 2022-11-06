package server

import (
	"Config/app/database"
	"Config/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type ConfigWrapper struct {
	proto.UnimplementedConfigWrapperServer
}

func (receiver *ConfigWrapper) GetConfig(ctx context.Context, in *proto.RequestService) (*proto.ResponseService, error) {
	conn := database.GetConnection().Pool

	getConfig := `
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
	
	err := conn.QueryRow(ctx, getConfig, in.GetService()).Scan(&result.Service, &result.Version, &result.Config, &result.ServiceID)
	if err != nil {
		return nil, err
	}

	log.Printf("Get config for %s", in.GetService())

	config := make(map[string]string, 0)
	err = json.Unmarshal(result.Config, &config)
	if err != nil{
		log.Println(err)
		return nil, err
	}

	return &proto.ResponseService{Service: result.Service, Version: int32(result.Version), Data: config}, nil
}
func (receiver *ConfigWrapper) CreateConfig(ctx context.Context, in *proto.RequestConfig) (*proto.Status, error) {
	conn := database.GetConnection().Pool

	firstVersion := int32(1)

	insertService := `
		INSERT INTO service 
		    (service, latest_version) 
		VALUES 
		       ($1, $2) 
		RETURNING service_id
	`

	var serviceID int
	err := conn.QueryRow(ctx, insertService, in.Service, firstVersion).Scan(&serviceID)
	if err != nil {
		return nil, err
	}

	insertConfig := `
		INSERT INTO config 
		    (service_id, version, cfg) 
		VALUES 
		       ($1, $2, $3) 
	`
	jsonCFG, err := json.Marshal(in.Data)
	if err != nil {
		return nil, err
	}
	fmt.Println(jsonCFG)
	conn.QueryRow(ctx, insertConfig, serviceID, firstVersion, jsonCFG)

	log.Printf("Create config for %s", in.Service)

	return &proto.Status{Message: "OK"}, nil
}
func (receiver *ConfigWrapper) DeleteUnusedConfig(ctx context.Context, in *proto.RequestService) (*proto.Status, error) {
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

	if latest_version > in.Version {
		DeleteConfig := `
			DELETE FROM config
			WHERE service_id = $1 AND version = $2
		`

		conn.QueryRow(ctx, DeleteConfig, service_id, in.Version)

		log.Printf("Config for service %s version %d has been deleted", in.Service, in.Version)

		return &proto.Status{Message: "Version has been deleted"}, nil
	}
	log.Printf("Config %s has one version or has no this version", in.Service)

	return &proto.Status{Message: "This version in used or has no this version"}, nil
}

func (receiver *ConfigWrapper) UpdateConfig(ctx context.Context, in *proto.RequestConfig) (*proto.Status, error) {
	conn := database.GetConnection().Pool

	GetServiceIdAndLatestVersion := `
		SELECT service_id, latest_version FROM service
		WHERE service = $1
	`
	var service_id, latest_version int32
	err := conn.QueryRow(ctx, GetServiceIdAndLatestVersion, in.Service).Scan(&service_id, &latest_version)
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

	jsonCFG, err := json.Marshal(in.Data)
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

	log.Printf("Update %s to %d version", in.Service, newVersion)

	return &proto.Status{Message: "OK"}, nil
}
