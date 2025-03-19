BEGIN TRANSACTION;

DROP INDEX IF EXISTS gauge_metric_id;
DROP INDEX IF EXISTS counter_metric_id;

DROP TABLE gauges;
DROP TABLE counters;

COMMIT; 