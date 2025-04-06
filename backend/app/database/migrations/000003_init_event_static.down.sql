DROP TABLE IF EXISTS event_statuses;

DELETE FROM event_formats WHERE format IN ('IN_PERSON', 'ONLINE', 'HYBRID');