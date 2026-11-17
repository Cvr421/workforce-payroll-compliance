CREATE TABLE compliance_obligations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    code text NOT NULL,
    name text NOT NULL,
    category text NOT NULL CHECK (category IN ('PF','ESI','PT','TDS','LABOUR','OTHER')),
    frequency text NOT NULL CHECK (frequency IN ('monthly','quarterly','annual')),
    due_day int NOT NULL CHECK (due_day BETWEEN 1 AND 28),
    jurisdiction text NOT NULL DEFAULT 'IN',
    source_ref text NOT NULL DEFAULT '',
    is_demo boolean NOT NULL DEFAULT true,
    UNIQUE (id, tenant_id),
    UNIQUE (tenant_id, code)
);

CREATE TABLE compliance_tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    establishment_id uuid,
    obligation_id uuid NOT NULL,
    period_label text NOT NULL,
    due_date date NOT NULL,
    status text NOT NULL DEFAULT 'not_started' CHECK (status IN
        ('not_started','in_progress','awaiting_evidence','under_review','completed')),
    priority text NOT NULL DEFAULT 'medium' CHECK (priority IN ('low','medium','high')),
    owner_id uuid REFERENCES users(id),
    reviewer_id uuid REFERENCES users(id),
    closed_by uuid REFERENCES users(id),
    closed_at timestamptz,
    reopen_count int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id, tenant_id) REFERENCES organizations (id, tenant_id),
    FOREIGN KEY (obligation_id, tenant_id) REFERENCES compliance_obligations (id, tenant_id),
    UNIQUE (id, tenant_id)
);
CREATE UNIQUE INDEX compliance_tasks_unique_period_uq
    ON compliance_tasks (organization_id, obligation_id, period_label,
                       COALESCE(establishment_id, '00000000-0000-0000-0000-000000000000'::uuid));
CREATE INDEX compliance_tasks_due_idx ON compliance_tasks (tenant_id, due_date) WHERE status <> 'completed';

CREATE TABLE compliance_evidence (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    task_id uuid NOT NULL,
    file_name text NOT NULL,
    content_type text NOT NULL,
    size_bytes bigint NOT NULL CHECK (size_bytes > 0),
    sha256 text NOT NULL,
    storage_path text NOT NULL,
    uploaded_by uuid NOT NULL REFERENCES users(id),
    review_status text NOT NULL DEFAULT 'pending' CHECK (review_status IN ('pending','accepted','rejected')),
    review_note text NOT NULL DEFAULT '',
    reviewed_by uuid REFERENCES users(id),
    reviewed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (task_id, tenant_id) REFERENCES compliance_tasks (id, tenant_id)
);

CREATE TABLE compliance_comments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    task_id uuid NOT NULL,
    author_id uuid NOT NULL REFERENCES users(id),
    body text NOT NULL CHECK (length(body) BETWEEN 1 AND 2000),
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (task_id, tenant_id) REFERENCES compliance_tasks (id, tenant_id)
);