package es

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uio"
	"github.com/sirupsen/logrus"
	"strings"
)

const BulkActionIndex = "index"
const BulkActionCreate = "create"

func BulkBodyFromLines(
	lineObjs []interface{},
) (body string, err error) {
	strLines := make([]string, 0)
	for _, line := range lineObjs {
		bLine, errM := json.Marshal(line)
		err = errM
		if nil != err {
			logrus.Errorf("json unmarshal bulk index line failed. line: %#v, error: %s.", line, err)
			return
		}

		strLines = append(strLines, string(bLine))
	}

	body = strings.Join(strLines, "\n") + "\n"
	return
}

// build bulk index request body from map[id]doc
func BulkBodyFromIdDocMap(
	action string,
	index string,
	docMap map[string]interface{},
) (body string, err error) {
	lineObjs := make([]interface{}, 0)
	for id, doc := range docMap {
		lineIndex := map[string]interface{}{
			action: map[string]interface{}{
				"_index": index,
				"_id":    id,
			},
		}

		lineObjs = append(lineObjs, lineIndex, doc)
	}

	body, err = BulkBodyFromLines(lineObjs)
	if nil != err {
		logrus.Errorf("get buck index body from lines failed. lines: %#v, error: %s.", lineObjs, err)
		return
	}

	return
}

func BulkBodyFromDocSlice(
	action string,
	index string,
	docSlice []interface{},
) (body string, err error) {
	lineObjs := make([]interface{}, 0)
	for _, doc := range docSlice {
		lineIndex := map[string]interface{}{
			action: map[string]interface{}{
				"_index": index,
			},
		}

		lineObjs = append(lineObjs, lineIndex, doc)
	}

	body, err = BulkBodyFromLines(lineObjs)
	if nil != err {
		logrus.Errorf("bulk body from lines failed. error: %s.", err)
		return
	}

	return
}

// handle response body
func HandleResponseBody(
	response *esapi.Response,
	outBody interface{},
) (err error) {
	bBody, err := uio.ReadAllAndClose(response.Body)
	if err != nil {
		logrus.Errorf("read response body failed. error: %s", err)
		return
	}
	if response.IsError() {
		err = RespError(bBody)
		if err != nil {
			logrus.Errorf("es response error. %s", err)
			return
		}
		return
	}

	if outBody == nil {
		return
	}

	err = json.Unmarshal(bBody, outBody)
	if err != nil {
		logrus.Errorf("unmarshal response body failed. %s", err)
		return
	}
	return
}
