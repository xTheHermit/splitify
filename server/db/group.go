package db

import "database/sql"

type GroupRepo interface {
	Select(int) (Group, error)
	Insert(string) (int, error)
	Delete(int) (bool, error)
	Update(Group) (bool, error)
}

type groupRepo struct {
	db *sql.DB
}

func newGroupRepo(db *sql.DB) (GroupRepo, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS user_group(" +
															"id INTEGER PRIMARY KEY AUTOINCREMENT," +
															"name TEXT);")
	if err != nil {
		return nil, err
	}
	return &groupRepo{
		db: db,
	}, nil
}

func (h *groupRepo) Select(id int) (Group, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans Group
	err := h.db.QueryRowContext(ctx, "SELECT id, name FROM user_group WHERE id = $1", id).
				Scan(&ans.Id, &ans.Name)
	if err != nil {
		return Group{}, err
	}
	return ans, nil
}

func (h *groupRepo) Insert(name string) (int, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var id int
	query := "INSERT INTO user_group(name) VALUES($1) RETURNING id"
	err := h.db.QueryRowContext(ctx, query, name).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (h *groupRepo) Delete(id int) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	_, err := h.db.ExecContext(ctx, "DELETE FROM user_group WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *groupRepo) Update(group Group) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	existing, err := h.Select(group.Id)
	if err != nil {
		return false, err
	}

	final := fillMerge(existing, group)

	_, err = h.db.ExecContext(ctx, "UPDATE user_group SET name = $2 WHERE id = $1", final.Id, final.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func fillMerge(source Group, extra Group) Group {
	if extra.Id > 0 {
		source.Id = extra.Id
	}
	if extra.Name != "" {
		source.Name = extra.Name
	}
	return source
}
