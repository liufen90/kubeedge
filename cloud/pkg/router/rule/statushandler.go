package rule

import (
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/cloud/pkg/common/modules"
	"github.com/kubeedge/kubeedge/cloud/pkg/router/constants"
	"github.com/kubeedge/kubeedge/cloud/pkg/router/messagelayer"
	"k8s.io/klog/v2"
	"time"
)

type ExecResult struct {
	RuleID    string
	ProjectID string
	Status    string
	Error     ErrorMsg
}

type ErrorMsg struct {
	Detail    string
	Timestamp time.Time
}

var ResultChannel chan ExecResult
var StopChan chan bool

func init() {
	StopChan = make(chan bool)
	go do(StopChan)
}

/*unc do(stop chan bool) {
	ResultChannel = make(chan ExecResult, 1024)
	r := <-ResultChannel
	if r.Status == "SUCCESS" {
		msg := model.NewMessage("")
		resource := "node/nodename/" + r.ProjectID + "/rule/" + r.RuleID
		msgv1 := msg.BuildRouter("router", "rule", resource, "update")
		msgv1.Content = r

		beehiveContext.Send("edgecontroller", *msgv1)
	}
}*/

func do(stop chan bool) {
	ResultChannel = make(chan ExecResult, 1024)

	for {
		select {
		case r := <-ResultChannel:
			msg := model.NewMessage("")
			resource, err := messagelayer.BuildResourceForRouter(r.ProjectID, model.ResourceTypeRuleStatus, r.RuleID) //message构建时需要的常数都在model里定义过
			if err != nil {
				klog.Warningf("build message resource failed with error: %s", err)
				continue
			}
			msg.Content = r
			msg.BuildRouter(modules.RouterModuleName, constants.GroupResource, resource, model.UpdateOperation) //modules里定义了
			beehiveContext.Send(modules.EdgeControllerModuleName, *msg)
			klog.V(4).Infof("send message successfully,operation: %s", msg.GetOperation(), msg.GetResource())
		case _, ok := <-stop:
			if !ok {
				klog.Warningf("do stop channel is closed")
			}
		}
		return
	}
}
