package drivers

import (
	"errors"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var ErrNoDrivers = errors.New("ms sql driver error: driver not found")
var ErrorInvalidDriver = errors.New("invalid driver")

const (
	NativeClient11 string = "SQL Server Native Client 11.0"
	ODBC13         string = "ODBC Driver 13 for SQL Server"
	SQLServer      string = "SQL Server"
)

var orderedDrivers = []string{
	SQLServer,
	NativeClient11,
	ODBC13,
}

func getDrivers() ([]string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\ODBC\ODBCINST.INI\ODBC Drivers`, registry.QUERY_VALUE)
	if err != nil {
		return nil, errors.New("openkey error")
	}
	defer k.Close()
	s, err := k.ReadValueNames(0)
	if err != nil {
		return nil, errors.New("read value error")
	}
	sort.Strings(s)
	return s, nil
}

func InstalledDrivers() ([]string, error) {
	var drivers []string
	d, err := getDrivers()
	if err != nil {
		return drivers, errors.New("getting drivers error")
	}
	for _, v := range d {
		for _, d := range orderedDrivers {
			if strings.EqualFold(d, v) {
				drivers = append(drivers, v)
			}
		}
	}
	return drivers, nil
}

func BestDriver() (string, error) {
	drivers, err := getDrivers()
	if err != nil {
		return "", errors.New("getting drivers error")
	}
	for _, d := range orderedDrivers {
		for _, v := range drivers {
			if strings.EqualFold(d, v) {
				return d, nil
			}
		}
	}
	return "", ErrNoDrivers
}

func ValidDriver(d string) error {
	drivers, err := InstalledDrivers()
	if err != nil {
		return errors.New("installed drivers error")
	}
	for _, v := range drivers {
		if v == d {
			return nil
		}
	}
	return ErrorInvalidDriver
}
