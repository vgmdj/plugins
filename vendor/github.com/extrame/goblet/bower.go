package goblet

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/extrame/go-bower/bower"
	"github.com/extrame/goblet/config"
)

var bower_cache = make(map[string][][2]string)

func (s *Server) Bower(name string, version ...string) (res [][2]string, err error) {
	if *s.env == config.ProductEnv {
		if res, ok := bower_cache[name]; ok {
			return res, nil
		}
	}

	root := filepath.Join(*s.wwwRoot, "public", "plugins", name)
	if _, err = os.Stat(root); os.IsNotExist(err) {
		if *s.env == config.ProductEnv {
			log.Panicf("no %s plugins in production environment", name)
		}
		if _, err = os.Stat(filepath.Join(*s.wwwRoot, "public", ".bowerrc")); os.IsNotExist(err) {
			ioutil.WriteFile(filepath.Join(*s.wwwRoot, "public", ".bowerrc"), []byte(`{"directory" : "plugins"}`), 0644)
		}
		if len(version) > 0 {
			name = name + "#" + version[0]
		}
		c := exec.Command("bower", "install", "-S", name, "--allow-root")
		c.Env = os.Environ()
		c.Dir = filepath.Join(*s.wwwRoot, "public")
		c.Stderr = LogFile
		if err = c.Run(); err != nil {
			return
		}
	}

	res = make([][2]string, 0)

	var bts []byte
	if bts, err = ioutil.ReadFile(filepath.Join(root, "bower.json")); err == nil {
		var b *bower.Component
		if b, err = bower.ParseBowerJSON(bts); err == nil {
			appendHTML(s, b, name, &res)
		}
	}

	if *s.env == config.ProductEnv {
		if err == nil {
			bower_cache[name] = res
		}
	}

	return
}

func appendHTML(s *Server, b *bower.Component, name string, maps *[][2]string) {

	root := filepath.Join(*s.wwwRoot, "public", "plugins")
	if b.Dependencies != nil {
		for k := range b.Dependencies {
			if bts, e := ioutil.ReadFile(filepath.Join(root, k, "bower.json")); e == nil {
				if b1, err := bower.ParseBowerJSON(bts); err == nil {
					appendHTML(s, b1, b1.Name, maps)
				} else {
					log.Println(err)
				}
			}
		}
	}

	res := ""
	switch bs := b.Main.(type) {
	case []interface{}:
		for _, v := range bs {
			res += appendHtmlItem(*s.env, root, name, v.(string))
		}
	case string:
		res += appendHtmlItem(*s.env, root, name, bs)
	default:
		log.Panicf("%v,%T", b.Main, b.Main)
	}

	*maps = append(*maps, [2]string{name, res})

	return
}

func appendHtmlItem(env, root, name, v string) string {
	if strings.HasSuffix(v, ".js") {
		if env == config.ProductEnv {
			//try to use min version
			min_v := strings.Replace(v, ".js", ".min.js", 1)
			if _, err := os.Stat(filepath.Join(root, name, min_v)); !os.IsNotExist(err) {
				return "<script src=/plugins/" + name + "/" + min_v + "></script>"
			}
		}
		return "<script src=/plugins/" + name + "/" + v + "></script>"
	} else if strings.HasSuffix(v, ".css") {
		return "<link href=/plugins/" + name + "/" + v + " rel='stylesheet'></link>"
	} else {
		return "<link href=/plugins/" + name + "/" + v + "></link>"
	}
}
