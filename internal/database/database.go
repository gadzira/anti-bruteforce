package database

import (
	"context"
	"fmt"
	"net"

	// import psql driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Entry struct {
	IP   string `json:"ip"`   // 192.1.1.0
	Mask string `json:"mask"` // 255.255.255.128
	List string `json:"list"` // black/white
}

type CIDR struct {
	cidr string
}

type DataBase struct {
	db *sqlx.DB
	l  *zap.Logger
}

func New(log *zap.Logger) DataBase {
	// nolint:exhaustivestruct
	return DataBase{
		l: log,
	}
}

func (s *DataBase) Connect(ctx context.Context, dsn string) error {
	var err error
	s.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	s.l.Info("Connected to database!")
	return s.db.PingContext(ctx)
}

func (s *DataBase) Close() error {
	return s.db.Close()
}

func (s *DataBase) AddToList(ctx context.Context, e *Entry) error {
	stringMask := net.IPMask(net.ParseIP(e.Mask).To4())
	length, _ := stringMask.Size()
	cidr := fmt.Sprintf("%s/%d", e.IP, length) // cidr = 192.1.1.0/25

	_, err := s.db.ExecContext(ctx, `INSERT INTO list (cidr, list) VALUES ($1, $2)`, cidr, e.List)
	if err != nil {
		return fmt.Errorf("cannot insert: %w", err)
	}
	return nil
}

func (s *DataBase) RemoveFromList(ctx context.Context, e *Entry) error {
	stringMask := net.IPMask(net.ParseIP(e.Mask).To4())
	length, _ := stringMask.Size()
	cidr := fmt.Sprintf("%s/%d", e.IP, length)

	_, err := s.db.ExecContext(ctx, `DELETE FROM list WHERE  cidr = $1 AND list = $2`, cidr, e.List)
	if err != nil {
		return fmt.Errorf("cannot delete: %w", err)
	}
	return nil
}

func (s *DataBase) CheckInList(ctx context.Context, ip string, list string) (bool, error) {
	//
	if err := s.db.Ping(); err != nil {
		s.db.Close()
		return false, err
	}
	//
	rows, err := s.db.QueryContext(ctx, `SELECT cidr FROM list WHERE list = $1`, list)
	if err != nil {
		return false, fmt.Errorf("can't check IP in %s list: %w", list, err)
	}
	defer rows.Close()

	var cidrs []CIDR
	for rows.Next() {
		var c CIDR
		if err := rows.Scan(&c.cidr); err != nil {
			return false, fmt.Errorf("can't scan rows: %w", err)
		}
		cidrs = append(cidrs, c)
	}

	for _, v := range cidrs {
		_, n, err := net.ParseCIDR(v.cidr)
		if err != nil {
			return false, fmt.Errorf("can't parse cidr: %w", err)
		}

		if n.Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}

	return false, nil
}
