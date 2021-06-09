package apps

import (
	"bufio"
	"bytes"
	"cloudiac/consts"
	"cloudiac/consts/e"
	"cloudiac/libs/ctx"
	"cloudiac/models"
	"cloudiac/services"
	"cloudiac/services/logstorage"
	"cloudiac/utils"
	"cloudiac/utils/logs"
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func FollowTaskLog(c *ctx.GinRequestCtx) error {
	taskId := c.Query("taskId")
	if taskId == "" {
		taskId = c.Query("taskGuid")
	}
	if taskId == "" {
		// logPath example: "logs/ct-c2j2g5rn8qhqp9ku9a6g/run-c2mdu4ecie6qs8gmsmkgg"
		logPath := c.Query("logPath")
		parts := strings.Split(logPath, "/")
		for _, s := range parts {
			if strings.HasPrefix(s, "run-") {
				taskId = s
				break
			}
		}
	}

	if taskId == "" {
		return e.New(e.BadRequest, http.StatusBadRequest, fmt.Errorf("'taskId' or 'logPath' required"))
	}

	logger := logs.Get().WithField("func", "FollowTaskLog").WithField("taskId", taskId)

	var (
		task *models.Task
		err  error
	)
	// 等待任务启动
	for i := 0; ; i++ {
		select {
		case <-c.Context.Done():
			break
		default:
		}

		task, err = services.GetTaskByGuid(c.ServiceCtx().DB(), taskId)
		if err != nil {
			if e.IsRecordNotFound(err) {
				return e.New(e.TaskNotExists)
			}
			logger.Errorf("query task err: %v", err)
			return e.New(e.DBError)
		}

		if !task.Started() {
			if i < 10 {
				time.Sleep(time.Second * time.Duration(i+1))
			} else {
				time.Sleep(time.Second * 10)
			}
			continue
		}
		break
	}

	var reader io.Reader
	if task.Exited() { // 己退出的任务直接读取全量日志
		if content, err := logstorage.Get().Read(task.BackendInfo.LogFile); err != nil {
			logger.Errorf("read task log error: %v", err)
			return err
		} else {
			reader = bytes.NewBuffer(content)
		}
	} else { // 否则实时从 runner 获取日志
		pr, pw := io.Pipe()
		reader = pr

		go func() {
			if err := fetchRunnerTaskLog(pw, task); err != nil {
				logger.Errorf("fetchRunnerTaskLog error: %v", err)
			}
		}()

		go func() {
			<-c.Request.Context().Done()
			logger.Tracef("connect closed")
			pr.Close()
		}()
	}

	scanner := bufio.NewScanner(reader)
	eventId := 0 // to indicate the message id
	for scanner.Scan() {
		c.Render(-1, sse.Event{
			Id:    strconv.Itoa(eventId),
			Event: "message",
			Data:  scanner.Text(),
		})
		c.Writer.Flush()
		eventId += 1
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// 从 runner 获取任务日志，直到任务结束
func fetchRunnerTaskLog(writer io.WriteCloser, task *models.Task) error {
	// close 后 read 端会触发 EOF error
	defer writer.Close()

	logger := logs.Get().WithField("func", "fetchRunnerTaskLog").WithField("taskId", task.Guid)

	runnerAddr := fmt.Sprintf("%v", task.BackendInfo.BackendUrl)

	params := url.Values{}
	params.Add("taskId", task.Guid)
	params.Add("templateId", task.TemplateGuid)
	params.Add("containerId", task.BackendInfo.ContainerId)
	wsConn, err := utils.WebsocketDail(runnerAddr, consts.RunnerTaskLogFollowURL, params)
	if err != nil {
		return err
	}
	defer func() {
		_ = utils.WebsocketClose(wsConn)
	}()

	for {
		_, reader, err := wsConn.NextReader()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				logger.Tracef("read message error: %v", err)
				return nil
			} else {
				logger.Errorf("read message error: %v", err)
				return err
			}
		} else {
			_, err := io.Copy(writer, reader)
			if err != nil {
				if err == io.ErrClosedPipe {
					return nil
				}
				logger.Infof("io.Copy: %v", err)
				return err
			}
		}
	}
}