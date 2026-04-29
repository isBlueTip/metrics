package repository

import (
	"context"
	"time"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	Pool *pgxpool.Pool
	Conn *pgx.Conn
}

func (s *DBStorage) SetGauge(name string, val float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "INSERT INTO GAUGES (NAME, VALUE)" +
		"VALUES ($1, $2)" +
		"ON CONFLICT (NAME) DO UPDATE " +
		"SET VALUE = EXCLUDED.VALUE;"

	_, err := s.Pool.Exec(ctx, sql, name, val)
	return err
}

func (s *DBStorage) SetCounter(name string, val int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "INSERT INTO COUNTERS (NAME, VALUE)" +
		"VALUES ($1, $2)" +
		"ON CONFLICT (NAME) DO UPDATE " +
		"SET VALUE = EXCLUDED.VALUE;"

	_, err := s.Pool.Exec(ctx, sql, name, val)
	return err
}

func (s *DBStorage) GetGauge(name string) (val float64, exists bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "SELECT name, value " +
		"FROM GAUGES " +
		"WHERE NAME = $1;"

	row := s.Pool.QueryRow(ctx, sql, name)

	var g models.GaugeModel
	err := row.Scan(&g.Name, &g.Value)
	if err != nil {
		return val, false
	}

	return g.Value, true
}

func (s *DBStorage) GetCounter(name string) (val int64, exists bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "SELECT name, value " +
		"FROM COUNTERS " +
		"WHERE NAME = $1;"

	row := s.Pool.QueryRow(ctx, sql, name)

	var c models.CounterModel
	err := row.Scan(&c.Name, &c.Value)
	if err != nil {
		return val, false
	}

	return c.Value, true
}

func (s *DBStorage) GetGauges() (res []models.GaugeModel, err error) {
	gauges := make([]models.GaugeModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "SELECT name, value " +
		"FROM GAUGES;"

	rows, err := s.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var g models.GaugeModel
		err = rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}
		gauges = append(gauges, g)
	}

	return gauges, nil
}

func (s *DBStorage) GetCounters() (res []models.CounterModel, err error) {
	gauges := make([]models.CounterModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "SELECT name, value " +
		"FROM COUNTERS;"

	rows, err := s.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c models.CounterModel
		err = rows.Scan(&c.Name, &c.Value)
		if err != nil {
			return nil, err
		}
		gauges = append(gauges, c)
	}

	return gauges, nil
}

func (s *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := s.Pool.Ping(ctx)

	return err
}

func NewDBStorage(connString string) (*DBStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	s := DBStorage{Pool: pool}

	err = s.initTables(ctx)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *DBStorage) initTables(ctx context.Context) error {
	gaugesSQL := "CREATE TABLE IF NOT EXISTS gauges(" +
		"\"name\" TEXT PRIMARY KEY," +
		"\"value\" FLOAT);"

	_, err := s.Pool.Exec(ctx, gaugesSQL)

	if err != nil {
		return err
	}

	countersSQL := "CREATE TABLE IF NOT EXISTS counters(" +
		"\"name\" TEXT PRIMARY KEY," +
		"\"value\" INTEGER);"

	_, err = s.Pool.Exec(ctx, countersSQL)

	return err
}
