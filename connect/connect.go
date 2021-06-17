package connect

import (
	"errors"
	"fmt"
	"github.com/svetlakoff88/dblib/drivers"
	"strings"
)

type Connection struct {
	driver              string
	Server              string
	User                string
	Password            string
	Trusted             bool
	Database            string
	MultiSubnetFailover bool
}

func (c *Connection) Driver(d string) string {
	return c.driver
}

func (c *Connection) SetDriver(d string) error {
	err := drivers.ValidDriver(d)
	if err != nil {
		return err
	}
	c.driver = d
	return nil
}

func (c *Connection) ConnectionString() (string, error) {
	var cxn string
	if c.driver == "" {
		driver, err := drivers.BestDriver()
		if err == drivers.ErrNoDrivers {
			return "", err
		}
		if err != nil {
			return "", errors.New("driver setting error: failed connection")
		}
		c.driver = driver
	}
	cxn += fmt.Sprintf("Driver={%s}; ", c.driver)
	if c.Server == "" {
		return "", errors.New("server error: server not found")
	}
	cxn += fmt.Sprintf("Server={%s}; ", c.Server)
	if c.Trusted || (c.User == "" && c.Password == "") {
		cxn += fmt.Sprintf("Trusted_Connection=yes; ")
	} else {
		cxn += fmt.Sprintf("UID=%s; PWD=%s; ", c.User, c.Password)
	}
	if c.Database != "" {
		cxn += fmt.Sprintf("Database=%s; ", c.Database)
	}
	if c.MultiSubnetFailover {
		cxn += "MultiSubnetFailover=Yes; "
	}
	cxn += strings.TrimSpace(cxn)
	return cxn, nil
}

func Parse(s string) (Connection, error) {
	var c Connection
	var err error
	s = strings.TrimSpace(s)
	attribs := strings.Split(strings.TrimSuffix(s, ";"), ";")
	for _, a := range attribs {
		p := strings.Split(a, "=")
		if len(p) != 2 {
			return c, errors.New("wrong attrib: " + a)
		}
		k := strings.ToLower(strings.TrimSpace(p[0]))
		v := strings.TrimSpace(p[1])
		v = strings.TrimPrefix(strings.TrimSuffix(v, "}"), "{")
		switch k {
		case "driver":
			err = c.SetDriver(v)
			if err != nil {
				return c, errors.New("set driver error")
			}
		case "server", "address", "addr":
			c.Server = v
		case "uid", "user id":
			c.User = v
		case "pwd", "password":
			c.Password = v
		case "database":
			c.Database = v
		case "trusted_connection":
			if strings.ToLower(v) == "yes" {
				c.Trusted = true
			}
		case "multisubnetfailover":
			if strings.ToLower(v) == "yes" {
				c.MultiSubnetFailover = true
			}
		}
	}
	return c, err
}
