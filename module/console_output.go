package module

import (
	"reflect"

	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化output插件映射关系表
	AddOutput("console", reflect.TypeOf(ConsoleOutput{}))
}

// Console output插件，将消息输出到控制台上
type ConsoleOutput struct {
	LifeCycle
	status StatusPlugin
}

func (c *ConsoleOutput) Send(msgs *Message) {
	if c.status != Started {
		log.Warnln("Console output is not running, output nothing.")
		return
	}

	// log.Printf("Output:\n\tHeader: %+v, Body: %+v\n", msgs.Data.Raw, msgs.Data.Target)
	// c.SetType("console").SetTarget("stdout").Send(msgs)

	// 全局动态变量
	var vars map[string]interface{}
	if vv, ok := msgs.Data.Items["vars"]; ok {
		vars = vv.(map[string]interface{})
	} else {
		vars = make(map[string]interface{})
	}
	// 获取hook配置 默认为console
	// 添加hooks不存在的默认配置console
	if sendtypes, ok := msgs.Data.Items["hooks"]; ok {
		if len(sendtypes.([]interface{})) > 0 {
			for _, types := range sendtypes.([]interface{}) {
				if t, ok := types.(map[interface{}]interface{})["type"]; ok {
					switch t.(string) {
					case "console":
						NewHookAdapter(nil).SetType("console").Send(msgs)
					case "email":
						email, err := NewEmail(types.(map[interface{}]interface{}), msgs, vars)
						if err != nil {
							log.WithFields(log.Fields{
								"plugin": "console_output",
								"type":   "email",
							}).Errorln(err)
						} else {
							NewHookAdapter(email).SetType("email").Send(msgs)
						}
					case "wechat":
						wechat, err := NewWechat(types.(map[interface{}]interface{}), msgs, vars)
						if err != nil {
							log.WithFields(log.Fields{
								"plugin": "console_output",
								"type":   "wechat",
							}).Errorln(err)
						} else {
							NewHookAdapter(wechat).SetType("wechat").Send(msgs)
						}
					default:
						log.Warnf("未适配该类型的hooks: %s", t.(string))
					}
				} else {
					log.Errorln("hooks 配置内容不包含[type]字段,请检查！")
					break
				}
			}
		}
	} else {
		NewHookAdapter(nil).SetType("console").Send(msgs)
	}

}

func (c *ConsoleOutput) Init() {
	c.name = "Console output"
	c.status = Started
}
