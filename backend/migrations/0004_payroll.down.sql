DROP TABLE IF EXISTS payslips;
DROP TABLE IF EXISTS payroll_events;
DROP TRIGGER IF EXISTS payroll_run_employees_immutable ON payroll_run_employees;
DROP FUNCTION IF EXISTS forbid_change_when_run_finalized();
DROP TABLE IF EXISTS payroll_run_employees;
DROP TABLE IF EXISTS payroll_runs;
DROP TABLE IF EXISTS statutory_rules;