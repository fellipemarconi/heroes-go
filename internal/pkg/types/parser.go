package types

import "github.com/jackc/pgx/v5/pgtype"

func ParseUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID

	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, err
	}

	return uuid, nil
}
