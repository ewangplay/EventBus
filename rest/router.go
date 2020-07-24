package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

type Router struct {
	opts   *config.EB_Options
	routes map[string]i.Handler
}

func NewRouter(opts *config.EB_Options, routes map[string]i.Handler) *Router {
	return &Router{opts, routes}
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//路由分发设置，用来判断url是否合法，通过配置文件的正则表达式配置

	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")

	resources, err := this.parseURL(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, this.makeErrorResult(-1, err.Error()))
		return
	}

	var res []string
	var resource string
	for i, r := range resources {
		if i == 0 {
			resource = r
		} else {
			resource += fmt.Sprintf("/%s", r)
		}
		res = append(res, resource)
	}

	var this_handler i.Handler
	var result string

	for i := len(res) - 1; i >= 0; i-- {
		resource := res[i]

		this_handler, err = this.getHandler(resource)
		if err == nil {
			result, err = this_handler.Process(r, resources, this_handler.ProcessFunc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, result)
			} else {
				header.Add("Content-Length", fmt.Sprintf("%v", len(result)))
				io.WriteString(w, result)
			}
			return
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, this.makeErrorResult(-1, err.Error()))
	}

	return
}

func (this *Router) getHandler(resource string) (i.Handler, error) {
	handler, found := this.routes[resource]
	if found && handler != nil {
		return handler, nil
	} else {
		return nil, errors.New("handler not found.")
	}
}

func (this *Router) parseURL(url string) (resources []string, err error) {
	//url pattern example: "/(v\\d+)/(\\w+)/?(\\w+)?"
	urlPattern := this.opts.EBUrlPattern
	urlRegexp, err := regexp.Compile(urlPattern)
	if err != nil {
		return
	}

	matchs := urlRegexp.FindStringSubmatch(url)
	if matchs == nil {
		err = errors.New("Wrong Request URL")
		return
	}

	for i := 1; i < len(matchs); i++ {
		resources = append(resources, matchs[i])
	}

	return
}

func (this *Router) makeErrorResult(errcode int, errmsg string) string {

	data := map[string]interface{}{
		c.ERROR_CODE: errcode,
		c.ERROR_MSG:  errmsg,
	}
	result, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("{\"%s\":%d,\"%s\":\"%s\"}", c.ERROR_CODE, errcode, c.ERROR_MSG, errmsg)
	}
	return string(result)
}
