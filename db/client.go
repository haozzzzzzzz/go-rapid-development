package db

import (
	"github.com/jmoiron/sqlx"
	"context"
	"github.com/sirupsen/logrus"
	"database/sql"
)

type ConnectionChecker interface {
	Before(client *Client) error
	After(err error) error
}

type CommandChecker interface {
	Before(strSql string, args ...interface{}) error
	After(err error) error
}

type Client struct {
	DB *sqlx.DB
	Ctx context.Context
	ConnectionCheckers []ConnectionChecker
	CommandCheckers []CommandChecker
}

func NewClient(ctx context.Context, sqlxDB * sqlx.DB) (client *Client, err error) {
	client = &Client{
		DB: sqlxDB,
		Ctx:ctx,
	}
	err = client.DB.PingContext(ctx)
	if nil != err {
		logrus.Errorf("ping db failed. %s", err)
		return
	}
	return
}

func (m *Client) QueryRowx(query string, args ...interface{}) (*sqlx.Row, error) {
	return m.DB.QueryRowxContext(m.Ctx, query, args...), nil
}

func (m *Client) Get(dest interface{}, query string, args ...interface{}) (err error) {
	err = m.DB.GetContext(m.Ctx, dest, query, args...)
	return
}

func (m *Client) NamedQuery(query string, arg interface{}) (rows *sqlx.Rows, err error) {
	rows, err = m.DB.NamedQueryContext(m.Ctx, query, arg)
	return
}

func (m *Client) Query(query string, args ...interface{})(*sqlx.Rows, error) {
	return m.DB.QueryxContext(m.Ctx, query, args...)
}

func (m *Client) Exec(query string, args ...interface{}) (sql.Result, error){
	return m.DB.ExecContext(m.Ctx, query, args...)
}

func (m *Client) NamedExec(query string, arg interface{}) (result sql.Result, err error) {
	return m.DB.NamedExecContext(m.Ctx, query, arg)
}

func (m *Client) Select(slice interface{}, query string, args ...interface{}) (err error) {
	return m.DB.SelectContext(m.Ctx, slice, query, args...)
}
