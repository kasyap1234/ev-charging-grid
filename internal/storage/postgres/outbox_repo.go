package postgres


type OutboxRepo  struct {
	pool *pgxpool.Pool
	queries db.Queries
}
