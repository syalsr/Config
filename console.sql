select service.service_id, service.service, config.version, config.cfg from service
INNER JOIN config ON service.service_id = config.service_id
WHERE "service" = 'managed-k8snsk1ss'
ORDER BY version DESC
LIMIT 1