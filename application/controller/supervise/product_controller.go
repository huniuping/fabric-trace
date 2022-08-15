package supervise

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	assemble_sdk "raft-fabric-project/application/sdk/assemble"
	contract_sdk "raft-fabric-project/application/sdk/contract"
	design_sdk "raft-fabric-project/application/sdk/design"
	manufacture_sdk "raft-fabric-project/application/sdk/manufacture"
	process_sdk "raft-fabric-project/application/sdk/process"
	quality_sdk "raft-fabric-project/application/sdk/quality"
	"strings"
)

func Trance(ctx *gin.Context) {
	//切片输出结果
	map_data:=map[string]interface{}{}
	gongdans:=[]map[string]interface{}{}
	//process_gongdans:=[]map[string]interface{}{}
	processes:=[]map[string]interface{}{}
	assembles:=[]map[string]interface{}{}
	qualifys:=[]map[string]interface{}{}
	drawings:=[]map[string]interface{}{}
	contracts:=[]map[string]interface{}{}
	//查询产品质检和组装
	product_id := ctx.Query("product_id")
	args := [][]byte{[]byte(product_id)}
	chaincode_name0 := "qualifycc"
	chaincode_name1 := "assemblecc"
	fnc := "query"
	rsp_quality,err :=quality_sdk.ChannelQuery(chaincode_name0,fnc,args)
	rsp_assemble,err :=assemble_sdk.ChannelQuery(chaincode_name1,fnc,args)
	if err != nil||string(rsp_quality.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询质检信息失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)
	if err != nil||string(rsp_assemble.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询组装信息失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)
	qualify := make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_quality.Payload)), &qualify)
	qualifys=append(qualifys,qualify)
	map_data["qualitys"]=qualifys
	assemble := make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_assemble.Payload)), &assemble)
	assembles=append(assembles,assemble)
	map_data["assembles"]=assembles

	//查询产品组件
	chaincode_name2 :="processcc"
	fnc1 :="queryByProductId"
	rsp_process,err :=process_sdk.ChannelQuery(chaincode_name2,fnc1,args)
	if err != nil||string(rsp_process.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询组件信息失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)

	payloads :=strings.Split(string(rsp_process.Payload),";;")
	for i := 0; i < len(payloads); i++ {
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		delete(temp_map,"sequential_aggregate_signature")
		processes=append(processes,temp_map)
	}
	map_data["processes"]=processes

	//查询组装工单
	assemble_gongdan_id :=assemble["work_order_id"].(string)
	fmt.Println(assemble_gongdan_id)
	args1 := [][]byte{[]byte(assemble_gongdan_id)}
	chaincode_name3 :="gongdancc"
	rsp_assemble_gongdan,err :=manufacture_sdk.ChannelQuery(chaincode_name3,fnc,args1)
	if err != nil||string(rsp_assemble_gongdan.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询组装工单失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)

	assemble_gongdan:=make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_assemble_gongdan.Payload)), &assemble_gongdan)
	gongdans=append(gongdans,assemble_gongdan)

	//组件工单
	process_gongdan_temp := make(map[string]interface{})//一行所有数据存在这一个map中
	json.Unmarshal([]byte(payloads[0]), &process_gongdan_temp)
	process_gongdan_id :=process_gongdan_temp["work_order_id"].(string)
	args2 := [][]byte{[]byte(process_gongdan_id)}
	rsp_process_gongdan,err :=manufacture_sdk.ChannelQuery(chaincode_name3,fnc,args2)
	if err != nil||string(rsp_process_gongdan.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询组件工单失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)
	process_gongdan:=make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_process_gongdan.Payload)), &process_gongdan)
	gongdans=append(gongdans,process_gongdan)
	map_data["gongdans"]=gongdans

	//查询合同
	contract_id :=process_gongdan["contract_id"].(string)
	args3 := [][]byte{[]byte(contract_id)}
	chaincode_name4 :="contractcc"
	rsp_contract,err :=contract_sdk.ChannelQuery(chaincode_name4,fnc,args3)
	if err != nil||string(rsp_contract.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询合同信息失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)
	contract:=make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_contract.Payload)), &contract)
	contracts=append(contracts,contract)
	map_data["contracts"]=contracts

	//查询图纸
	drawing_id :=process_gongdan["drawing_id"].(string)
	args4 := [][]byte{[]byte(drawing_id)}
	chaincode_name5 :="drawingcc"
	rsp_drawing,err :=design_sdk.ChannelQuery(chaincode_name5,fnc,args4)
	if err != nil||string(rsp_drawing.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "查询图纸信息失败",
			"data": nil,
		})
		return
	}
	fmt.Println(err)
	drawing:=make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp_drawing.Payload)), &drawing)
	drawings=append(drawings,drawing)
	map_data["drawings"]=drawings


	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}


