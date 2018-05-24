package db

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CommandChecker interface {
	Before(client *Client, strSql string, args ...interface{})
	After(err error)
}

type CommandCheckerMaker interface {
	NewChecker() CommandChecker
}

type Client struct {
	DB                  *sqlx.DB
	Ctx                 context.Context
	CommandCheckerMaker CommandCheckerMaker
}

func (m *Client) CommandChecker() CommandChecker {
	if m.CommandCheckerMaker == nil {
		return nil
	}

	return m.CommandCheckerMaker.NewChecker()
}

func NewClient(ctx context.Context, sqlxDB *sqlx.DB, commCheckerMaker CommandCheckerMaker) (client *Client, err error) {
	client = &Client{
		DB:                  sqlxDB,
		Ctx:                 ctx,
		CommandCheckerMaker: commCheckerMaker,
	}

	return
}

func (m *Client) Ping() (err error) {
	err = m.DB.PingContext(m.Ctx)
	if nil != err {
		logrus.Errorf("ping db failed. %s", err)
		return
	}

	return
}

func (m *Client) QueryRowx(query string, args ...interface{}) (row *sqlx.Row, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, args...)
		defer func() {
			checker.After(err)
		}()
	}

	row = m.DB.QueryRowxContext(m.Ctx, query, args...)
	return
}

func (m *Client) Get(dest interface{}, query string, args ...interface{}) (err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, args...)
		defer func() {
			checker.After(err)
		}()
	}

	err = m.DB.GetContext(m.Ctx, dest, query, args...)
	return
}

func (m *Client) NamedQuery(query string, arg interface{}) (rows *sqlx.Rows, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, arg)
		defer func() {
			checker.After(err)
		}()
	}
	rows, err = m.DB.NamedQueryContext(m.Ctx, query, arg)
	return
}

func (m *Client) Query(query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, args...)
		defer func() {
			checker.After(err)
		}()
	}

	rows, err = m.DB.QueryxContext(m.Ctx, query, args...)
	return
}

func (m *Client) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, args...)
		defer func() {
			checker.After(err)
		}()
	}
	result, err = m.DB.ExecContext(m.Ctx, query, args...)
	return
}

func (m *Client) NamedExec(query string, arg interface{}) (result sql.Result, err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, arg)
		defer func() {
			checker.After(err)
		}()
	}
	result, err = m.DB.NamedExecContext(m.Ctx, query, arg)
	return
}

func (m *Client) Select(slice interface{}, query string, args ...interface{}) (err error) {
	// checker
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, query, args...)
		defer func() {
			checker.After(err)
		}()
	}
	err = m.DB.SelectContext(m.Ctx, slice, query, args...)
	return
}
