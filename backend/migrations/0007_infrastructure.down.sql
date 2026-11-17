DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS notifications;
DROP TRIGGER IF EXISTS payroll_events_append_only ON payroll_events;
DROP TRIGGER IF EXISTS audit_events_append_only ON audit_events;
DROP FUNCTION IF EXISTS forbid_update_delete();
DROP TABLE IF EXISTS audit_events;