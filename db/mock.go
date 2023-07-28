package db

import "github.com/pashagolub/pgxmock/v2"

func MockConnection() (SqlDB, pgxmock.PgxPoolIface, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		return nil, nil, err
	}

	pc := mock.(SqlDB)
	sqlDbConn = pc
	return pc, mock, nil
}
