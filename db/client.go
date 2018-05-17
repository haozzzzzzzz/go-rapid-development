package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ConnectionChecker interface {
	Before(client *Client)
	After(err error)
}

type ConnectionCheckerMaker interface {
	NewChecker() ConnectionChecker
}

type CommandChecker interface {
	Before(strSql string, args ...interface{})
	After(err error)
}

type CommandCheckerMaker interface {
	NewChecker() CommandChecker
}

type Client struct {
	DB                     *sqlx.DB
	Ctx                    context.Context
	ConnectionCheckerMaker ConnectionCheckerMaker
	CommandCheckerMaker    CommandCheckerMaker
}

func (m *Client) ConnectionChecker() ConnectionChecker {
	if m.ConnectionCheckerMaker == nil {
		return nil
	}
	return m.ConnectionCheckerMaker.NewChecker()
}

func (m *Client) CommandChecker() CommandChecker {
	if m.CommandCheckerMaker == nil {
		return nil
	}

	return m.CommandCheckerMaker.NewChecker()
}

func NewClient(ctx context.Context, sqlxDB *sqlx.DB) (client *Client, err error) {
	client = &Client{
		DB:  sqlxDB,
		Ctx: ctx,
	}

	// checker
	checker := client.ConnectionChecker()
	if checker != nil {
		checker.Before(client)
		defer checker.After(err)
	}

	err = client.DB.PingContext(ctx)
	if nil != err {
		logrus.Errorf("ping db failed. %s", err)
		return
	}
	return
}

func (m *Client) QueryRowx(query string, args ...interface{}) (rows *sqlx.Row, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, args...)
		defer checker.After(err)
	}

	rows = m.DB.QueryRowxContext(m.Ctx, query, args...)
	return
}

func (m *Client) Get(dest interface{}, query string, args ...interface{}) (err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, args...)
		defer checker.After(err)
	}

	err = m.DB.GetContext(m.Ctx, dest, query, args...)
	return
}

func (m *Client) NamedQuery(query string, arg interface{}) (rows *sqlx.Rows, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, arg)
		defer checker.After(err)
	}
	rows, err = m.DB.NamedQueryContext(m.Ctx, query, arg)
	return
}

func (m *Client) Query(query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, args...)
		defer checker.After(err)
	}

	rows, err = m.DB.QueryxContext(m.Ctx, query, args...)
	return
}

func (m *Client) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, args...)
		defer checker.After(err)
	}
	result, err = m.DB.ExecContext(m.Ctx, query, args...)
	return
}

func (m *Client) NamedExec(query string, arg interface{}) (result sql.Result, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, arg)
		defer checker.After(err)
	}
	result, err = m.DB.NamedExecContext(m.Ctx, query, arg)
	return
}

func (m *Client) Select(slice interface{}, query string, args ...interface{}) (err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(query, args...)
		defer checker.After(err)
	}
	err = m.DB.SelectContext(m.Ctx, slice, query, args...)
	return
}
