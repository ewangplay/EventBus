package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

type BaseHandler struct {
	Opts *config.EB_Options
	i.ILogger
	i.IIDCounter
	i.IProducer
	i.IJobManager
}

func NewBaseHandler(opts *config.EB_Options, logger i.ILogger, idcounter i.IIDCounter, messager i.IProducer, jobmgr i.IJobManager) *BaseHandler {
	return &BaseHandler{opts, logger, idcounter, messager, jobmgr}
}

func (this *BaseHandler) Process(r *http.Request, resources []string, f i.ProcessFunc) (string, error) {
	var err error
	var body []byte
	var startTime time.Time

	startTime = time.Now()

	result := make(map[string]interface{})

	//Generate the log id
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	log_id1 := rand.Intn(100000)
	log_id2 := rand.Intn(100000)
	log_id := fmt.Sprintf("%d%d", log_id1, log_id2)
	result[c.LOG_ID] = log_id

	this.Info("[LOG_ID:%v] [METHOD:%v] [URL:%v]", log_id, r.Method, r.RequestURI)

	//Parse the request parameters
	params, err := this.parseArgs(r)
	if err != nil {
		result[c.ERROR_CODE] = -1
		result[c.ERROR_MSG] = "parse request parameter error" //err.Error()
		goto END
	}

	//Read the request body
	body, err = ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		result[c.ERROR_CODE] = -1
		result[c.ERROR_MSG] = "read request body error" //err.Error()
		goto END
	}

	//Perform the actual business process
	err = f(r.Method, resources, params, body, result)
	if err != nil {
		result[c.ERROR_CODE] = -1
		if strings.HasPrefix(err.Error(), "[ERROR_INFO]") {
			result[c.ERROR_MSG] = strings.TrimPrefix(err.Error(), "[ERROR_INFO]")
		} else {
			result[c.ERROR_MSG] = fmt.Sprintf("systerm error! LOG_ID: %v", log_id)
		}
		goto END
	}

	result[c.ERROR_CODE] = 0

END:
	if err != nil {
		this.Error("[LOG_ID:%v] %v", log_id, err)
		if string(body) != "" {
			this.Error("[LOG_ID:%v] [Request Body : %v]", log_id, string(body))
		}
		this.Error("[LOG_ID:%v] [Response Result : %v]", log_id, result)
	}

	result[c.REQUEST_URL] = r.RequestURI
	result[c.TIME_COST] = fmt.Sprintf("%v", time.Since(startTime))
	this.Info("[LOG_ID:%v] [COST:%v]", log_id, result[c.TIME_COST])

	resStr, _ := this.createJSON(result)

	return resStr, err
}

func (this *BaseHandler) createJSON(result map[string]interface{}) (string, error) {
	r, err := json.Marshal(result)
	if err != nil {
		this.Error("%v", err)
		return "", err
	}
	return string(r), nil
}

func (this *BaseHandler) parseArgs(r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	//每次都重新生成一个新的map，否则之前请求的参数会保留其中
	res := make(map[string]string)
	for k, v := range r.Form {
		res[k] = v[0]
	}

	return res, nil
}
