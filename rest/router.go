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

// Router struct define
type Router struct {
	opts   *config.EBOptions
	routes map[string]i.Handler
}

// NewRouter ...
func NewRouter(opts *config.EBOptions, routes map[string]i.Handler) *Router {
	return &Router{opts, routes}
}

// ServerHTTP ...
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//路由分发设置，用来判断url是否合法，通过配置文件的正则表达式配置

	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")

	resources, err := router.parseURL(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, router.makeErrorResult(-1, err.Error()))
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

	var handler i.Handler
	var result string

	for i := len(res) - 1; i >= 0; i-- {
		resource := res[i]

		handler, err = router.getHandler(resource)
		if err == nil {
			result, err = handler.Process(r, resources, handler.ProcessFunc)
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
		io.WriteString(w, router.makeErrorResult(-1, err.Error()))
	}

	return
}

func (router *Router) getHandler(resource string) (i.Handler, error) {
	handler, found := router.routes[resource]
	if found && handler != nil {
		return handler, nil
	}
	return nil, errors.New("handler not found")
}

func (router *Router) parseURL(url string) (resources []string, err error) {
	//url pattern example: "/(v\\d+)/(\\w+)/?(\\w+)?"
	urlPattern := router.opts.EBUrlPattern
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

func (router *Router) makeErrorResult(errcode int, errmsg string) string {

	data := map[string]interface{}{
		c.ErrorCode: errcode,
		c.ErrorMsg:  errmsg,
	}
	result, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("{\"%s\":%d,\"%s\":\"%s\"}", c.ErrorCode, errcode, c.ErrorMsg, errmsg)
	}
	return string(result)
}
