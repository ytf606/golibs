package db

import (
	"errors"
)

type Conn interface {
	Build() (err error)
}

// DBConn ...
type DBConn struct {
	name string
	conf *Config
	// DB default
	db *DB
}

var _ Conn = (*DBConn)(nil)

func New() *DBConn {
	return &DBConn{}
}

func (c *DBConn) SetName(name string) *DBConn {
	c.name = name
	return c
}

func (c *DBConn) SetConf(config *Config) *DBConn {
	c.conf = config
	return c
}

// Build ...
func (c *DBConn) Build() (err error) {
	if c.name == "" && c.conf == nil {
		if c.conf.OnDialError == "panic" {
			panic(err)
		}
		return errors.New("please set conn name or config")
	}

	// writer
	_db, err := Open(c.conf)
	if err != nil {
		if c.conf.OnDialError == "panic" {
			panic(err)
		}
		return
	}
	sqlWDB, err := _db.DB()
	if err != nil {
		return
	}
	if err = sqlWDB.Ping(); err != nil {
		return
	}

	c.db = _db
	instances.Store(c.name, _db)
	return nil
}

func (c *DBConn) GetDB() *DB {
	return c.db
}

func (c *DBConn) Close() error {
	_db, _ := c.db.DB()
	err := _db.Close()
	return err
}

// ClusterConn ...
type ClusterConn struct {
	name string
	conf *ClusterConfig
	// w Writer
	w *DB
	// R Reader
	r *DB
}

var _ Conn = (*ClusterConn)(nil)

func NewCluster() *ClusterConn {
	return &ClusterConn{}
}

func (c *ClusterConn) SetName(name string) *ClusterConn {
	c.name = name
	return c
}

func (c *ClusterConn) SetConf(config *ClusterConfig) *ClusterConn {
	c.conf = config
	return c
}

// Build cluster ...
func (c *ClusterConn) Build() (err error) {
	if c.name == "" && c.conf == nil {
		err = errors.New("please set conn name or config")
		panic(err)
	}

	// writer
	_wdb, err := Open(c.conf.W)
	if err != nil {
		if c.conf.W.OnDialError == "panic" {
			panic(err)
		}
		return
	}
	sqlWDB, err := _wdb.DB()
	if err != nil {
		return
	}
	if err = sqlWDB.Ping(); err != nil {
		return
	}

	// reader
	_rdb, err := Open(c.conf.R)
	if err != nil {
		if c.conf.R.OnDialError == "panic" {
			panic(err)
		}
		return
	}
	sqlRDB, err := _rdb.DB()
	if err != nil {
		return
	}
	if err = sqlRDB.Ping(); err != nil {
		return
	}

	c.w = _wdb
	c.r = _rdb
	clusterInstances.Store(c.name, c)
	return nil
}

func (c *ClusterConn) GetW() *DB {
	return c.w
}

func (c *ClusterConn) GetR() *DB {
	return c.r
}

func (c *ClusterConn) Close() error {
	_db, _ := c.w.DB()
	err := _db.Close()
	if err != nil {
		return err
	}
	_db, _ = c.r.DB()
	err = _db.Close()
	return err
}
