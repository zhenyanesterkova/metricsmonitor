BEGIN TRANSACTION;

CREATE TABLE gauges(
    id VARCHAR(200) UNIQUE NOT NULL PRIMARY KEY,
    g_value DOUBLE PRECISION NOT NULL
);

CREATE TABLE counters(
    id VARCHAR(200) UNIQUE NOT NULL PRIMARY KEY,
    delta BIGINT NOT NULL
);

CREATE INDEX gauge_metric_id ON gauges (id);
CREATE INDEX counter_metric_id ON counters (id);

COMMIT;