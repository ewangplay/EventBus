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

// BaseHandler structure
type BaseHandler struct {
	Opts *config.EBOptions
	i.Logger
	i.Counter
	i.Producer
	i.JobManager
}

// NewBaseHandler new instance
func NewBaseHandler(opts *config.EBOptions, logger i.Logger, idcounter i.Counter, messager i.Producer, jobmgr i.JobManager) *BaseHandler {
	return &BaseHandler{opts, logger, idcounter, messager, jobmgr}
}

// Process the http Request
func (bh *BaseHandler) Process(r *http.Request, resources []string, f i.ProcessFunc) (string, error) {
	var err error
	var body []byte
	var startTime time.Time

	startTime = time.Now()

	result := make(map[string]interface{})

	//Generate the log id
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	logID1 := rand.Intn(100000)
	logID2 := rand.Intn(100000)
	logID := fmt.Sprintf("%d%d", logID1, logID2)
	result[c.LogID] = logID

	bh.Info("[LogID:%v] [METHOD:%v] [URL:%v]", logID, r.Method, r.RequestURI)

	//Parse the request parameters
	params, err := bh.parseArgs(r)
	if err != nil {
		result[c.ErrorCode] = -1
		result[c.ErrorMsg] = "parse request parameter error" //err.Error()
		goto END
	}

	//Read the request body
	body, err = ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		result[c.ErrorCode] = -1
		result[c.ErrorMsg] = "read request body error" //err.Error()
		goto END
	}

	//Perform the actual business process
	err = f(r.Method, resources, params, body, result)
	if err != nil {
		result[c.ErrorCode] = -1
		if strings.HasPrefix(err.Error(), "[ERROR_INFO]") {
			result[c.ErrorMsg] = strings.TrimPrefix(err.Error(), "[ERROR_INFO]")
		} else {
			result[c.ErrorMsg] = fmt.Sprintf("systerm error! LogID: %v", logID)
		}
		goto END
	}

	result[c.ErrorCode] = 0

END:
	if err != nil {
		bh.Error("[LogID:%v] %v", logID, err)
		if string(body) != "" {
			bh.Error("[LogID:%v] [Request Body : %v]", logID, string(body))
		}
		bh.Error("[LogID:%v] [Response Result : %v]", logID, result)
	}

	result[c.RequestURL] = r.RequestURI
	result[c.TimeCost] = fmt.Sprintf("%v", time.Since(startTime))
	bh.Info("[LogID:%v] [COST:%v]", logID, result[c.TimeCost])

	resStr, _ := bh.createJSON(result)

	return resStr, err
}

func (bh *BaseHandler) createJSON(result map[string]interface{}) (string, error) {
	r, err := json.Marshal(result)
	if err != nil {
		bh.Error("%v", err)
		return "", err
	}
	return string(r), nil
}

func (bh *BaseHandler) parseArgs(r *http.Request) (map[string]string, error) {
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
