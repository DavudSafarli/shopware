package postgres

import (
	"context"
	"database/sql"
	"redirectware/internal"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// AddFullMatchRule upserts a redirect rule for an exact path.
func (s *Storage) AddFullMatchRule(ctx context.Context, rule *internal.FullMatchRule) error {
	const q = `INSERT INTO redirect_rules_full_match(from_raw, from_canonical, target)
              VALUES ($1, $2, $3)
              ON CONFLICT (from_canonical) DO UPDATE SET
                from_raw = EXCLUDED.from_raw,
                target = EXCLUDED.target`
	_, err := s.db.ExecContext(ctx, q, rule.FromRaw, rule.FromCanonical, rule.Target)
	return err
}

func (s *Storage) FindFullMatchRule(ctx context.Context, canonicalPath string) (rule internal.FullMatchRule, ok bool, err error) {
	const q = `SELECT from_raw, from_canonical, target FROM redirect_rules_full_match WHERE from_canonical = $1 LIMIT 1`
	var r internal.FullMatchRule
	err = s.db.QueryRowContext(ctx, q, canonicalPath).Scan(&r.FromRaw, &r.FromCanonical, &r.Target)
	if err == sql.ErrNoRows {
		return r, false, nil
	}
	if err != nil {
		return r, false, err
	}
	return r, true, nil
}

// GetWelcomePageURL returns the configured welcome page URL.
func (s *Storage) GetWelcomePageURL(ctx context.Context) (string, error) {
	const q = `SELECT target FROM welcome_page_url ORDER BY id ASC LIMIT 1`
	var url string
	err := s.db.QueryRowContext(ctx, q).Scan(&url)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return url, nil
}

// SetWelcomePageURL upserts the welcome page URL (single row semantics).
func (s *Storage) SetWelcomePageURL(ctx context.Context, url string) error {
	// delete all and insert one.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if _, err = tx.ExecContext(ctx, `TRUNCATE TABLE welcome_page_url`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO welcome_page_url(target) VALUES ($1)`, url); err != nil {
		return err
	}
	return tx.Commit()
}
