package db

import "database/sql"

type UserRepo interface {
	Select(int) (User, error)
	SelectByUsername(string) (User, error)
	Insert(string, string) (int, error)
	Update(int, string) (bool, error)
}

type userRepo struct {
	db *sql.DB
}

func newUserRepo(db *sql.DB) (UserRepo, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS user(id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		return nil, err
	}
	return &userRepo{
		db: db,
	}, nil
}

func (h *userRepo) Select(id int) (User, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans User
	query := "SELECT id, username, password FROM user WHERE id = $1"
	err := h.db.QueryRowContext(ctx, query, id).Scan(&ans.Id, &ans.Username, &ans.Password)
	if err != nil {
		return User{}, err
	}
	return ans, nil
}

func (h *userRepo) SelectByUsername(username string) (User, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans User
	query := "SELECT id, username, password FROM user WHERE username = $1"
	err := h.db.QueryRowContext(ctx, query, username).Scan(&ans.Id, &ans.Username, &ans.Password)
	if err != nil {
		return User{}, err
	}
	return ans, nil
}

func (h *userRepo) Insert(username string, password string) (int, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var id int
	query := "INSERT INTO user(username, password) VALUES($1, $2) RETURNING id"
	err := h.db.QueryRowContext(ctx, query, username, password).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (h *userRepo) Update(id int, password string) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	query := "UPDATE user SET password = $2 WHERE id = $1"
	_, err := h.db.ExecContext(ctx, query, id, password)

	if err != nil {
		return false, err
	}

	return true, nil
}
