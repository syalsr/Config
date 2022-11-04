SELECT 
		service.service, config.version, config.cfg, service.service_id 
	FROM service
	INNER JOIN 
		config ON service.service_id = config.service_id
		AND
		config.version = (SELECT MAX(config.version) FROM config
			WHERE config.service_id = (SELECT service_id FROM service
					WHERE "service" = $1))