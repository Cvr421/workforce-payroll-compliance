CREATE TABLE audit_events (
    id bigserial PRIMARY KEY,
    tenant_id uuid NOT NULL,
    actor_id uuid,
    action text NOT NULL,
    resource_type text NOT NULL,
    resource_id uuid,
    metadata jsonb NOT NULL DEFAULT '{}',
    request_id text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_events_resource_idx ON audit_events (tenant_id, resource_type, resource_id, created_at DESC);

CREATE FUNCTION forbid_update_delete() RETURNS trigger AS $$
BEGIN RAISE EXCEPTION '% on % is not allowed (append-only table)', TG_OP, TG_TABLE_NAME; END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER audit_events_append_only BEFORE UPDATE OR DELETE ON audit_events
    FOR EACH ROW EXECUTE FUNCTION forbid_update_delete();
CREATE TRIGGER payroll_events_append_only BEFORE UPDATE OR DELETE ON payroll_events
    FOR EACH ROW EXECUTE FUNCTION forbid_update_delete();

CREATE TABLE notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    kind text NOT NULL,
    title text NOT NULL,
    body text NOT NULL DEFAULT '',
    link text NOT NULL DEFAULT '',
    dedupe_key text NOT NULL,
    read_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, dedupe_key)
);
CREATE INDEX notifications_user_idx ON notifications (user_id, read_at, created_at DESC);

CREATE TABLE jobs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid,
    kind text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    status text NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','done','failed')),
    attempts int NOT NULL DEFAULT 0,
    max_attempts int NOT NULL DEFAULT 5,
    run_after timestamptz NOT NULL DEFAULT now(),
    locked_at timestamptz,
    locked_by text,
    last_error text NOT NULL DEFAULT '',
    dedupe_key text UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz
);
CREATE INDEX jobs_claim_idx ON jobs (run_after, id) WHERE status = 'queued';