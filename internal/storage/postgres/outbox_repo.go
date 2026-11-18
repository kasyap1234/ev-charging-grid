package postgres


type OutboxRepo  struct {
	pool *pgxpool.Pool
	queries db.Queries
}


var _domain.OutboxRepo=(*OutboxRepo)(nil)


func NewOutboxRepo(pool*pgxpool)*OutboxRepo{
	return &OutboxRepo{
		pool: pool,
		queries: db.New(pool)
	}
}
