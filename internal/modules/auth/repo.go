package auth

import "github.com/jmoiron/sqlx"

type Repo struct{ db *sqlx.DB }

func NewRepo(db *sqlx.DB) *Repo { return &Repo{db: db} }

type userRow struct {
	UserID       int64  `db:"user_id"`
	NIP          string `db:"nip"`
	FullName     string `db:"full_name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	RoleCode     string `db:"role_code"`
	RegionalID   int64  `db:"regional_id"`
	IsActive     bool   `db:"is_active"`
}

func (r *Repo) GetUserByNIP(nip string) (*userRow, error) {
	q := `
SELECT 
  u.user_id, u.nip, u.full_name, u.email, u.password_hash,
  r.code AS role_code,
  COALESCE(u.regional_id, 0) AS regional_id,
  COALESCE(u.is_active, true) AS is_active
FROM users u
JOIN roles r ON r.role_id = u.role_id
WHERE u.nip = $1
LIMIT 1;
`
	var row userRow
	if err := r.db.Get(&row, q, nip); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetUserByID(userID int64) (*userRow, error) {
	q := `
SELECT 
  u.user_id, u.nip, u.full_name, u.email, u.password_hash,
  r.code AS role_code,
  COALESCE(u.regional_id, 0) AS regional_id,
  COALESCE(u.is_active, true) AS is_active
FROM users u
JOIN roles r ON r.role_id = u.role_id
WHERE u.user_id = $1
LIMIT 1;
`
	var row userRow
	if err := r.db.Get(&row, q, userID); err != nil {
		return nil, err
	}
	return &row, nil
}
