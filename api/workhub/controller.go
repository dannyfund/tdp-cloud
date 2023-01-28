package workhub

import (
	"errors"

	"github.com/gin-gonic/gin"

	"tdp-cloud/module/dborm/user"
	"tdp-cloud/module/workhub"
)

// 节点列表

func list(c *gin.Context) {

	userId := c.GetUint("UserId")

	res := workhub.NodesOfUser(userId)

	c.Set("Payload", res)

}

// 执行脚本

type execParam struct {
	WorkerId string
	Payload  workhub.ExecPayload
}

func exec(c *gin.Context) {

	var rq *execParam

	if c.ShouldBind(&rq) != nil {
		c.Set("Error", "请求参数错误")
		return
	}

	send := workhub.NewSender(rq.WorkerId)

	if send == nil {
		c.Set("Error", "客户端已断开连接")
		return
	}

	if id, err := send.Exec(&rq.Payload); err == nil {
		c.Set("Payload", map[string]any{
			"Message": "命令下发完成",
			"TaskId":  id,
		})
	} else {
		c.Set("Error", err)
	}

}

// 注册节点

func register(c *gin.Context) {

	u, err := user.Fetch(&user.FetchParam{
		AppId: c.Param("appid"),
	})

	if err != nil || u.Id == 0 {
		c.AbortWithError(400, errors.New("授权失败"))
		return
	}

	c.Set("UserId", u.Id)
	c.Set("MachineId", c.Param("mid"))

	if err := workhub.Register(c); err != nil {
		c.AbortWithError(500, err)
		return
	}

}
