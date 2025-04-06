CREATE TABLE IF NOT EXISTS event_statuses (
status VARCHAR PRIMARY KEY
);

INSERT INTO event_statuses(status) VALUES
('ACTIVE'),
('DONE'),
('POSTPONED'),
('CANCELLED'),
('TRANSFERRED');

INSERT INTO event_formats(format) VALUES
('IN_PERSON'),
('ONLINE'),
('HYBRID');

