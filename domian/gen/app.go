package gen

import (
	"go2gen/conf"
	"log"
	"os"
	"reflect"
)

func (m *Manager) App(rootPath string, name string) error {
	tmpl := conf.Global.GetRealTmpl(rootPath, name)
	m.Tmpl = tmpl
	rv := reflect.ValueOf(tmpl)
	for i := 0; i < rv.Elem().NumField(); i++ {
		rs := rv.Elem().Field(i).Interface()
		dir := rs.(string)
		if dir != "" {
			if err := os.MkdirAll(dir, 0766); err != nil {
				log.Printf("Mkdir[%s] err: %v \n", dir, err)
			}
		}
	}
	return nil
}
