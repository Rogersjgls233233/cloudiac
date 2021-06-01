package runner

import (
	"bytes"
	"cloudiac/configs"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
)

type ContainerStatus struct {
	Status          *types.ContainerState
	LogContent      []string
}

func Run(req *http.Request) (string, error) {
	c, state, iacTemplate, err := ReqToCommand(req)
	if err != nil {
		return "", err
	}

	// TODO(ZhengYue):
	// 1. 根据模板ID创建目录(目录创建规则：/{template_uuid}/{task_id}/)，用于保存日志文件及挂载provider、密钥等文件
	// 2. 若需要保存模板状态，则根据参数生成状态配置文件，放入模板目录中，挂载至容器内部

	templateUUID := iacTemplate.TemplateUUID
	taskId := iacTemplate.TaskId
	templateDir, err := CreateTemplatePath(templateUUID, taskId)
	if err != nil {
		return "", err
	}
	// if state.SaveState != false {
	GenStateFile(state.StateBackendAddress, state.Scheme, state.StateKey, templateDir, state.SaveState)
	// }
	err = c.Create(templateDir)
	return c.ContainerInstance.ID, err
}

func Cancel(req *http.Request) error {
	task, err := ReqToTask(req)
	if err != nil {
		return err
	}
	err = task.Cancel()
	return err
}

type TaskLogsResp struct {
	LogContent      []string
	LogContentLines int
}

func GetTaskLogs(req *http.Request) (TaskLogsResp, error) {
	conf := configs.Get()
	task, err := ReqToTask(req)
	if err != nil {
		return TaskLogsResp{}, err
	}
	tlr := TaskLogsResp{}
	templateDir := fmt.Sprintf("%s/%s/%s", conf.Runner.LogBasePath, task.TemplateId, task.TaskId)
	logFile := fmt.Sprintf("%s/%s", templateDir, ContainerLogFileName)
	file, err := ioutil.ReadFile(logFile)
	if err != nil {
		return TaskLogsResp{}, err
	}
	buf := bytes.NewBuffer(file)
	MaxLines, _ := LineCounter(buf)

	lines, err := ReadLogFile(logFile, task.LogContentOffset, MaxLines)

	//logContent, err := FetchTaskLog(task.TemplateId, task.TaskId, task.LogContentOffset)
	if err != nil {
		return tlr, err
	}
	tlr.LogContentLines = len(lines)
	tlr.LogContent = lines
	return tlr, nil
}
