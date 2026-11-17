CREATE TABLE audits (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    establishment_id uuid,
    name text NOT NULL,
    audit_type text NOT NULL CHECK (audit_type IN ('labour','vendor','establishment')),
    period_start date NOT NULL,
    period_end date NOT NULL,
    auditor_id uuid NOT NULL REFERENCES users(id),
    status text NOT NULL DEFAULT 'planned' CHECK (status IN ('planned','in_progress','completed')),
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    UNIQUE (id, tenant_id),
    CHECK (period_end >= period_start)
);

CREATE TABLE audit_findings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    audit_id uuid NOT NULL,
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    severity text NOT NULL CHECK (severity IN ('critical','high','medium','low')),
    owner_id uuid NOT NULL REFERENCES users(id),
    target_date date NOT NULL,
    status text NOT NULL DEFAULT 'open' CHECK (status IN ('open','action_in_progress','pending_verification','closed')),
    closed_by uuid REFERENCES users(id),
    closed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (audit_id, tenant_id) REFERENCES audits (id, tenant_id),
    UNIQUE (id, tenant_id)
);

CREATE TABLE corrective_actions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    finding_id uuid NOT NULL,
    description text NOT NULL,
    owner_id uuid NOT NULL REFERENCES users(id),
    due_on date NOT NULL,
    status text NOT NULL DEFAULT 'open' CHECK (status IN ('open','evidence_submitted','verified','closed')),
    evidence_ref text NOT NULL DEFAULT '',
    reject_note text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (finding_id, tenant_id) REFERENCES audit_findings (id, tenant_id),
    CHECK (status NOT IN ('evidence_submitted','verified','closed') OR evidence_ref <> '')
);
CREATE INDEX corrective_actions_due_idx ON corrective_actions (tenant_id, due_on) WHERE status <> 'closed';