package db

import "database/sql"

type GroupMembershipRepo interface {
	SelectByGroup(int) (GroupMembership, error)
	SelectByUser(int) (GroupMembership, error)
	SelectByUserGroup(int, int) (GroupMembership, error)
	Insert(int, int, int) (bool, error)
	Delete(int, int) (bool, error)
	Update(int, int, int) (bool, error)
}

type groupMembershipRepo struct {
	db *sql.DB
}

func NewGroupMembershipRepo(db *sql.DB) (GroupMembershipRepo, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	query := "CREATE TABLE IF NOT EXISTS membership(" +
										"group_id INTEGER NOT NULL," +
										"user_id INTEGER NOT NULL," +
										"balance INTEGER," +
										"PRIMARY KEY (group_id, user_id));"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &groupMembershipRepo{
		db: db,
	}, nil
}

func (h *groupMembershipRepo) SelectByGroup(groupId int) (GroupMembership, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans GroupMembership
	query := "SELECT group_id, user_id, balance FROM membership WHERE group_id = $1"
	err := h.db.QueryRowContext(ctx, query, groupId).
			Scan(&ans.GroupId, &ans.UserId, &ans.Balance)
	if err != nil {
		return GroupMembership{}, err
	}
	return ans, nil
}

func (h *groupMembershipRepo) SelectByUser(userId int) (GroupMembership, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans GroupMembership
	query := "SELECT group_id, user_id, balance FROM membership WHERE user_id = $1"
	err := h.db.QueryRowContext(ctx, query, userId).
			Scan(&ans.GroupId, &ans.UserId, &ans.Balance)
	if err != nil {
		return GroupMembership{}, err
	}
	return ans, nil
}

func (h *groupMembershipRepo) SelectByUserGroup(userId int, groupId int) (GroupMembership, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	var ans GroupMembership
	query := "SELECT group_id, user_id, balance FROM membership WHERE user_id = $1 AND group_id = $2"
	err := h.db.QueryRowContext(ctx, query, userId, groupId).
			Scan(&ans.GroupId, &ans.UserId, &ans.Balance)
	if err != nil {
		return GroupMembership{}, err
	}
	return ans, nil
}

func (h *groupMembershipRepo) Insert(userId int, groupId int, balance int) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	query := "INSERT INTO membership(user_id, group_id, balance) VALUES($1, $2, $3)"
	_, err := h.db.ExecContext(ctx, query, userId, groupId, balance)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *groupMembershipRepo) Delete(userId int, groupId int) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	query := "DELETE FROM membership WHERE user_id = $1 AND group_id = $2"
	_, err := h.db.ExecContext(ctx, query, userId, groupId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *groupMembershipRepo) Update(userId int, groupId int, balance int) (bool, error) {
	ctx, cancel := getStandardContext()
	defer cancel()

	query := "UPDATE membership SET balance = $1 WHERE user_id = $2 AND group_id = $3"
	_, err := h.db.ExecContext(ctx, query, balance, userId, groupId)
	if err != nil {
		return false, err
	}

	return true, nil
}


