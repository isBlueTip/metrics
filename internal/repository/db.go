package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code[:2] == "08"
	}

	return false
}

type DBStorage struct {
	Pool *pgxpool.Pool
	Conn *pgx.Conn
}

func (s *DBStorage) SetGauge(name string, val float64) error {
	var err error
	for i := 0; i < 4; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		sql := "INSERT INTO GAUGES (NAME, VALUE)" +
			"VALUES ($1, $2)" +
			"ON CONFLICT (NAME) DO UPDATE " +
			"SET VALUE = EXCLUDED.VALUE;"

		_, err = s.Pool.Exec(ctx, sql, name, val)
		cancel()

		if err == nil || !isRetryableError(err) {
			return err
		}
		if i < 3 {
			time.Sleep(time.Duration(i*2+1) * time.Second)
		}
	}
	return err
}

func (s *DBStorage) SetCounter(name string, val int64) error {
	var err error
	for i := 0; i < 4; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		sql := "INSERT INTO COUNTERS (NAME, VALUE)" +
			"VALUES ($1, $2)" +
			"ON CONFLICT (NAME) DO UPDATE " +
			"SET " +
			"VALUE = COUNTERS.VALUE + EXCLUDED.VALUE;"

		_, err = s.Pool.Exec(ctx, sql, name, val)
		cancel()

		if err == nil || !isRetryableError(err) {
			return err
		}
		if i < 3 {
			time.Sleep(time.Duration(i*2+1) * time.Second)
		}
	}
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
	counters := make([]models.CounterModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := "SELECT name, value FROM COUNTERS;"

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
		counters = append(counters, c)
	}

	return counters, nil
}

func (s *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := s.Pool.Ping(ctx)

	return err
}

func (s *DBStorage) SaveToFile(path string) error {
	gauges, err := s.GetGauges()
	if err != nil {
		return err
	}

	counters, err := s.GetCounters()
	if err != nil {
		return err
	}

	gaugeMap := make(map[string]float64)
	for _, g := range gauges {
		gaugeMap[g.Name] = g.Value
	}

	counterMap := make(map[string]int64)
	for _, c := range counters {
		counterMap[c.Name] = c.Value
	}

	data := FileData{
		Gauge:   gaugeMap,
		Counter: counterMap,
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func (s *DBStorage) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var data FileData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for k, v := range data.Gauge {
		sql := "INSERT INTO GAUGES (NAME, VALUE) " +
			"VALUES ($1, $2) " +
			"ON CONFLICT (NAME) DO UPDATE " +
			"SET " +
			"VALUE = EXCLUDED.VALUE;"
		_, err := s.Pool.Exec(ctx, sql, k, v)
		if err != nil {
			return err
		}
	}

	for k, v := range data.Counter {
		sql := "INSERT INTO COUNTERS (NAME, VALUE) " +
			"VALUES ($1, $2) " +
			"ON CONFLICT (NAME) DO UPDATE " +
			"SET " +
			"VALUE = COUNTERS.VALUE + EXCLUDED.VALUE;"
		_, err := s.Pool.Exec(ctx, sql, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DBStorage) UpdateBatch(metrics []models.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range metrics {
		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				sql := "INSERT INTO GAUGES (NAME, VALUE) " +
					"VALUES ($1, $2) " +
					"ON CONFLICT (NAME) DO UPDATE " +
					"SET " +
					"VALUE = EXCLUDED.VALUE;"
				_, err = tx.Exec(ctx, sql, m.ID, *m.Value)
				if err != nil {
					return err
				}
			}
		case models.Counter:
			if m.Delta != nil {
				sql := "INSERT INTO COUNTERS (NAME, VALUE) " +
					"VALUES ($1, $2) " +
					"ON CONFLICT (NAME) DO UPDATE " +
					"SET " +
					"VALUE = COUNTERS.VALUE + EXCLUDED.VALUE;"
				_, err = tx.Exec(ctx, sql, m.ID, *m.Delta)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit(ctx)
}

func NewDBStorage(connString string) (*DBStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	s := DBStorage{Pool: pool, Conn: conn}

	err = s.initTables(ctx)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *DBStorage) initTables(ctx context.Context) error {
	gaugesSQL := `CREATE TABLE IF NOT EXISTS gauges(` +
		`"name" TEXT PRIMARY KEY,` +
		`"value" FLOAT);`

	_, err := s.Pool.Exec(ctx, gaugesSQL)

	if err != nil {
		return err
	}

	countersSQL := `CREATE TABLE IF NOT EXISTS counters(` +
		`"name" TEXT PRIMARY KEY,` +
		`"value" BIGINT);`

	_, err = s.Pool.Exec(ctx, countersSQL)

	return err
}
