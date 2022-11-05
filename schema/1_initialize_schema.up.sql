CREATE TABLE IF NOT EXISTS service
(
    service_id SERIAL PRIMARY KEY,
    service VARCHAR(100) NOT NULL UNIQUE,
    latest_version integer NOT NULL
);

CREATE TABLE IF NOT EXISTS config
(
    service_id INT REFERENCES service(service_id),
    version INT NOT NULL,
    cfg JSON NOT NULL
);