package outbox

import (
	"context"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository interface {
	Save(ctx context.Context, messages ...*outbox.Message) error
	GetNotPublishedMessages(ctx context.Context) ([]*outbox.Message, error)
}

var _ OutboxRepository = &repository{}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) (OutboxRepository, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &repository{
		db: db,
	}, nil
}

func (r *repository) Save(ctx context.Context, messages ...*outbox.Message) error {
	query := `insert into outbox(id, type, content, occurred, processed)
			  values ($1, $2, $3, $4, $5)
			  on conflict (id)
				 do update set processed = EXCLUDED.processed;`

	for _, m := range messages {
		_, err := r.db.Exec(ctx, query, m.ID, m.Name, m.Payload, m.OccurredAtUtc, m.ProcessedAtUtc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetNotPublishedMessages(ctx context.Context) ([]*outbox.Message, error) {

	empty := make([]*outbox.Message, 0)

	query := `select id, type, content, occurred
			  from outbox
			  where processed is null
			  order by occurred
			  limit 100`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return empty, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	messages := make([]*outbox.Message, 0, 100)
	for rows.Next() {

		m := outbox.Message{}
		err = rows.Scan(&m.ID, &m.Name, &m.Payload, &m.OccurredAtUtc)
		if err != nil {
			return empty, err
		}

		messages = append(messages, &m)
	}

	if err = rows.Err(); err != nil {
		return empty, err
	}

	return messages, nil
}
