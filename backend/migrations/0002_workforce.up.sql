CREATE TABLE employees (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    establishment_id uuid NOT NULL,
    employee_code text NOT NULL,
    full_name text NOT NULL,
    email text NOT NULL DEFAULT '',
    phone text NOT NULL DEFAULT '',
    designation text NOT NULL DEFAULT '',
    department text NOT NULL DEFAULT '',
    date_of_joining date NOT NULL,
    date_of_exit date,
    exit_reason text,
    status text NOT NULL DEFAULT 'active' CHECK (status IN ('active','exited')),
    bank_account_last4 text NOT NULL DEFAULT '' CHECK (bank_account_last4 ~ '^[0-9]{0,4}$'),
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    FOREIGN KEY (establishment_id, tenant_id) REFERENCES establishments (id, tenant_id),
    UNIQUE (id, tenant_id),
    UNIQUE (organization_id, employee_code),
    CHECK (date_of_exit IS NULL OR date_of_exit >= date_of_joining),
    CHECK ((status = 'exited') = (date_of_exit IS NOT NULL))
);
CREATE INDEX employees_search_idx ON employees (tenant_id, organization_id, status, full_name);

CREATE TABLE salary_structures (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    employee_id uuid NOT NULL,
    effective_from date NOT NULL,
    effective_to date,
    basic_paise bigint NOT NULL CHECK (basic_paise > 0),
    hra_paise bigint NOT NULL DEFAULT 0 CHECK (hra_paise >= 0),
    conveyance_paise bigint NOT NULL DEFAULT 0 CHECK (conveyance_paise >= 0),
    special_allowance_paise bigint NOT NULL DEFAULT 0 CHECK (special_allowance_paise >= 0),
    created_by uuid NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (employee_id, tenant_id) REFERENCES employees (id, tenant_id),
    CHECK (effective_to IS NULL OR effective_to > effective_from),
    EXCLUDE USING gist (employee_id WITH =, daterange(effective_from, effective_to, '[)') WITH &&)
);