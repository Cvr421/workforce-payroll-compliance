CREATE TABLE statutory_rules (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    rule_code text NOT NULL,
    description text NOT NULL,
    jurisdiction text NOT NULL DEFAULT 'IN',
    effective_from date NOT NULL,
    effective_to date,
    params jsonb NOT NULL,
    source_ref text NOT NULL DEFAULT '',
    version int NOT NULL DEFAULT 1,
    approval_status text NOT NULL DEFAULT 'draft' CHECK (approval_status IN ('draft','approved')),
    is_demo boolean NOT NULL DEFAULT true,
    EXCLUDE USING gist (tenant_id WITH =, rule_code WITH =, daterange(effective_from, effective_to, '[)') WITH &&)
);

CREATE TABLE payroll_runs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    period_year int NOT NULL CHECK (period_year BETWEEN 2000 AND 2100),
    period_month int NOT NULL CHECK (period_month BETWEEN 1 AND 12),
    status text NOT NULL DEFAULT 'draft' CHECK (status IN
        ('draft','calculated','awaiting_approval','approved','rejected','finalized','cancelled')),
    created_by uuid NOT NULL REFERENCES users(id),
    submitted_by uuid REFERENCES users(id),
    approved_by uuid REFERENCES users(id),
    rejected_by uuid REFERENCES users(id),
    rejected_reason text,
    finalized_by uuid REFERENCES users(id),
    finalized_at timestamptz,
    calc_version text NOT NULL DEFAULT 'v1',
    rules_snapshot jsonb,
    total_gross_paise bigint NOT NULL DEFAULT 0,
    total_net_paise bigint NOT NULL DEFAULT 0,
    error_count int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    UNIQUE (id, tenant_id)
);
CREATE UNIQUE INDEX payroll_runs_one_live_uq
    ON payroll_runs (organization_id, period_year, period_month)
    WHERE status NOT IN ('rejected','cancelled');
CREATE INDEX payroll_runs_list_idx ON payroll_runs (tenant_id, organization_id, status, created_at DESC);

CREATE TABLE payroll_run_employees (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    run_id uuid NOT NULL,
    employee_id uuid NOT NULL,
    employee_code text NOT NULL,
    employee_name text NOT NULL,
    days_in_month int NOT NULL,
    paid_half_days int NOT NULL,
    structure_snapshot jsonb NOT NULL,
    earned_gross_paise bigint NOT NULL,
    pf_wages_paise bigint NOT NULL,
    employee_pf_paise bigint NOT NULL,
    employer_pf_paise bigint NOT NULL,
    employee_esi_paise bigint NOT NULL,
    employer_esi_paise bigint NOT NULL,
    net_paise bigint NOT NULL,
    issues jsonb NOT NULL DEFAULT '[]',
    FOREIGN KEY (run_id, tenant_id) REFERENCES payroll_runs (id, tenant_id),
    UNIQUE (run_id, employee_id)
);

CREATE FUNCTION forbid_change_when_run_finalized() RETURNS trigger AS $$
DECLARE s text;
BEGIN
    SELECT status INTO s FROM payroll_runs WHERE id = OLD.run_id;
    IF s = 'finalized' THEN
        RAISE EXCEPTION 'payroll run % is finalized and immutable', OLD.run_id USING ERRCODE = '23000';
    END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
    RETURN NEW;
END; $$ LANGUAGE plpgsql;
CREATE TRIGGER payroll_run_employees_immutable
    BEFORE UPDATE OR DELETE ON payroll_run_employees
    FOR EACH ROW EXECUTE FUNCTION forbid_change_when_run_finalized();

CREATE TABLE payroll_events (
    id bigserial PRIMARY KEY,
    tenant_id uuid NOT NULL,
    run_id uuid NOT NULL,
    from_status text NOT NULL,
    to_status text NOT NULL,
    actor_id uuid NOT NULL REFERENCES users(id),
    reason text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (run_id, tenant_id) REFERENCES payroll_runs (id, tenant_id)
);

CREATE TABLE payslips (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    run_id uuid NOT NULL,
    employee_id uuid NOT NULL,
    payslip_number text NOT NULL,
    data jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (run_id, tenant_id) REFERENCES payroll_runs (id, tenant_id),
    UNIQUE (run_id, employee_id),
    UNIQUE (tenant_id, payslip_number)
);