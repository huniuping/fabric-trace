package all_router

import (
	"github.com/gin-gonic/gin"
	"raft-fabric-project/application/controller/assemble"
	"raft-fabric-project/application/controller/design"
	"raft-fabric-project/application/controller/process"
	"raft-fabric-project/application/controller/quality"
	"raft-fabric-project/application/controller/supervise"
	"raft-fabric-project/application/controller/contract"
	"raft-fabric-project/application1/controller/workOrder"
	"raft-fabric-project/application1/controller/login"
)

func InitRouter(router *gin.Engine) {

	assemble_group := router.Group("/assemble")
	design_group := router.Group("/design")
	process_group := router.Group("/process")
	quality_group := router.Group("/quality")
	contract_group := router.Group("/contract")
	supervise_group := router.Group("/supervise")
	workOrder_group := router.Group("/workerOrder")
	login_group := router.Group("/login")



	assemble.Router(assemble_group)
	design.Router(design_group)
	process.Router(process_group)
	supervise.Router(supervise_group)
	quality.Router(quality_group)
	contract.Router(contract_group)
	workOrder.Router(workOrder_group)
	login.Router(login_group)

}
