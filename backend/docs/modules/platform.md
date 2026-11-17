Problem: The platform foundation owns configuration, database access, HTTP envelopes, and process middleware.
Invariant: Requests are bounded and traceable, while database calls use contexts and a bounded pool.
Enforcement: The API applies request IDs, recovery, logging, security headers, CORS, body limits, and pgx pooling.
Proof: `go test -race ./...`, `go vet ./...`, and the `/readyz` smoke check cover this foundation.
