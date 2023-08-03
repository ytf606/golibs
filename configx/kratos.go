package configx

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var (
	g_kratos config.Config
)

func InitKratos(path string, bc interface{}) error {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)

	if err := c.Load(); err != nil {
		return err
	}

	if err := c.Scan(bc); err != nil {
		return err
	}
	g_kratos = c
	return nil
}

func GetValue(key string) string {
	ret, _ := g_kratos.Value(key).String()
	return ret
}

func GetKratos() config.Config {
	return g_kratos
}

func Close() {
	g_kratos.Close()
}
