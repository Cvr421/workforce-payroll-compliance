CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE tenants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    email text NOT NULL,
    full_name text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL CHECK (role IN ('admin','hr_admin','payroll_manager','compliance_officer','auditor')),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_uq ON users (lower(email));
CREATE INDEX users_tenant_idx ON users (tenant_id);

CREATE TABLE organizations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    name text NOT NULL,
    employer_model text NOT NULL CHECK (employer_model IN ('outsourcing','payroll_processing')),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (id, tenant_id),
    UNIQUE (tenant_id, name)
);

CREATE TABLE establishments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    name text NOT NULL,
    state text NOT NULL,
    address text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    UNIQUE (id, tenant_id),
    UNIQUE (organization_id, name)
);
CREATE INDEX establishments_org_idx ON establishments (tenant_id, organization_id);