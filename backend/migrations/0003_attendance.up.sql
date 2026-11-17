CREATE TABLE attendance_records (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    employee_id uuid NOT NULL,
    work_date date NOT NULL,
    status text NOT NULL CHECK (status IN ('present','absent','leave','weekly_off','holiday','half_day')),
    overtime_minutes int NOT NULL DEFAULT 0 CHECK (overtime_minutes >= 0),
    note text NOT NULL DEFAULT '',
    updated_by uuid NOT NULL REFERENCES users(id),
    updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (employee_id, tenant_id) REFERENCES employees (id, tenant_id),
    UNIQUE (employee_id, work_date)
);
CREATE INDEX attendance_date_idx ON attendance_records (tenant_id, work_date);

CREATE TABLE attendance_locks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    year int NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    month int NOT NULL CHECK (month BETWEEN 1 AND 12),
    locked_by uuid NOT NULL REFERENCES users(id),
    locked_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    UNIQUE (organization_id, year, month)
);