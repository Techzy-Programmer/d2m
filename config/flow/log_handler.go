package flow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
)

type deployLogHandler struct {
	msg   string
	title string
	steps [][]string
}

func (d *deployLogHandler) log(steps []string) *deployLogHandler {
	d.steps = append(d.steps, steps)
	return d
}

func (d *deployLogHandler) logStep(step string) *deployLogHandler {
	infoLvlStr := strconv.FormatUint(uint64(infoLvl), 10)
	d.log([]string{infoLvlStr, step})
	return d
}

func (d *deployLogHandler) logOk(steps string) *deployLogHandler {
	okLvlStr := strconv.FormatUint(uint64(okLvl), 10)
	d.log([]string{okLvlStr, steps})
	return d
}

func (d *deployLogHandler) logWarn(steps string) *deployLogHandler {
	warnLvlStr := strconv.FormatUint(uint64(warnLvl), 10)
	d.log([]string{warnLvlStr, steps})
	return d
}

func (d *deployLogHandler) logErr(steps string) *deployLogHandler {
	errLvlStr := strconv.FormatUint(uint64(errLvl), 10)
	d.log([]string{errLvlStr, steps})
	return d
}

func (d *deployLogHandler) save(lvl uint) {
	stepsJSON, err := json.Marshal(d.steps)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	deployLogStorage = append(deployLogStorage, db.DeploymentLog{
		Timestamp: time.Now().Unix(),
		DeployID:  currentDeployID,
		Steps:     string(stepsJSON),
		Title:     d.title,
		Message:   d.msg,
		Level:     lvl,
	})
}

func (d *deployLogHandler) reset(title string, msg string) {
	d.steps = [][]string{}
	d.title = title
	d.msg = msg
}

func newDeployLogHandler(title string, msg string) *deployLogHandler {
	logHandler := &deployLogHandler{}
	logHandler.reset(title, msg)
	return logHandler
}
