package repositories

import (
	"context"
	"fmt"
	"jwt-service/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) CreateSession(session *entities.Session) error {
	query := `INSERT INTO sessions (id, user_guid, refresh_hash, jti, ip_address, expires_at, used)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(context.Background(), query,
		session.ID,
		session.UserGUID,
		session.RefreshHash,
		session.JTI,
		session.IPAddress,
		session.ExpiresAt,
		session.Used)
	return err
}

func (r *SessionRepo) FindSessionByJTI(jti string) (*entities.Session, error) {
	query := `SELECT id, user_guid, refresh_hash, jti, ip_address, expires_at, used, created_at
              FROM sessions WHERE jti = $1`
	row := r.db.QueryRow(context.Background(), query, jti)

	var session entities.Session
	err := row.Scan(
		&session.ID,
		&session.UserGUID,
		&session.RefreshHash,
		&session.JTI,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.Used,
		&session.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}

	return &session, err
}

func (r *SessionRepo) MarkSessionAsUsed(id string) error {
	query := `UPDATE sessions SET used = true WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
