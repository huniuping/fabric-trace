package process

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	process_sdk "raft-fabric-project/application/sdk/process"
	"strings"
)

func Query(ctx *gin.Context) {
	process_product_id := ctx.Query("process_product_id")

	chaincode_name := "processcc"
	fnc := "query"
	args := [][]byte{[]byte(process_product_id)}
	rsp, err := process_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	fmt.Println("=============")
	xx:=string(rsp.Payload)

	fmt.Println(xx)
	map_data := make(map[string]interface{})
	json.Unmarshal([]byte(string(rsp.Payload)), &map_data)

	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func QueryByProductId(ctx *gin.Context) {
	producted_id := ctx.Query("producted_id")

	chaincode_name := "processcc"
	fnc := "queryByProductId"
	args := [][]byte{[]byte(producted_id)}
	rsp, err := process_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	payload_s:=string(rsp.Payload)
	payloads :=strings.Split(payload_s,";;")

	map_data:=map[string]interface{}{}
	processs:=[]map[string]interface{}{}

	for i := 0; i < len(payloads); i++ {
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		processs=append(processs,temp_map)
	}
	map_data["processs"]=processs
	fmt.Println("2222222222222222222222")
	fmt.Println(map_data)
	fmt.Println("222222222222222222222")
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func QueryByworkOrderId(ctx *gin.Context) {
	work_order_Id := ctx.Query("work_order_Id")
	chaincode_name := "processcc"
	fnc := "queryByContractId"
	args := [][]byte{[]byte(work_order_Id)}
	rsp, err := process_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	payload_s:=string(rsp.Payload)
	payloads :=strings.Split(payload_s,";;")

	map_data:=map[string]interface{}{}
	work_order:=[]map[string]interface{}{}

	for i := 0; i < len(payloads); i++ {
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		work_order=append(work_order,temp_map)
	}
	map_data["work_order_Id"]=work_order_Id
	fmt.Println(map_data)

	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,

	})
	return
}

func QueryProcessDetailByID(ctx *gin.Context) {
	process_product_id := ctx.Query("process_product_id")

	chaincode_name := "processcc"
	fnc := "query"
	args := [][]byte{[]byte(process_product_id)}
	rsp, err := process_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	payload_s:=string(rsp.Payload)
	payloads :=strings.Split(payload_s,";;")

	map_data:=map[string]interface{}{}
	processs:=[]map[string]interface{}{}

	for i := 0; i < len(payloads); i++ {
		if payloads[i] ==""{
			continue
		}
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		processs=append(processs,temp_map)
	}
	map_data["processsHistory"]=processs
	fmt.Println("2222222222222222222222")
	fmt.Println(map_data)
	fmt.Println("222222222222222222222")
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func Set(ctx *gin.Context) {
	process_product_id := ctx.PostForm("process_product_id")
	product_name := ctx.PostForm("product_name")
	work_order_id := ctx.PostForm("work_order_id")
	producted_id := ctx.PostForm("producted_id")
	technology := ctx.PostForm("producted_id")
	technology_sequence := ctx.PostForm("technology_sequence")
	sequential_aggregate_signature := ctx.PostForm("sequential_aggregate_signature")
	chaincode_name := "processcc"
	fnc := "set"
	args := [][]byte{[]byte(process_product_id),
		[]byte(product_name),
		[]byte(work_order_id),
		[]byte(producted_id),
		[]byte(technology),
		[]byte(technology_sequence),
		[]byte(sequential_aggregate_signature)}
	_, err := process_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经添加成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:9051]: Chaincode status Code: (500) UNKNOWN. Description: 该产品标识号已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该半成品标识号已存在",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "添加失败",
		})
		return
	}
}

func Update(ctx *gin.Context) {
	process_product_id := ctx.PostForm("process_product_id")
	product_name := ctx.PostForm("product_name")
	work_order_id := ctx.PostForm("work_order_id")
	producted_id := ctx.PostForm("producted_id")
	technology := ctx.PostForm("producted_id")
	technology_sequence := ctx.PostForm("technology_sequence")
	sequential_aggregate_signature := ctx.PostForm("sequential_aggregate_signature")
	chaincode_name := "processcc"
	fnc := "update"
	args := [][]byte{[]byte(process_product_id),
		[]byte(product_name),
		[]byte(work_order_id),
		[]byte(producted_id),
		[]byte(technology),
		[]byte(technology_sequence),
		[]byte(sequential_aggregate_signature)}
	_, err := process_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经更新成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:9051]: Chaincode status Code: (500) UNKNOWN. Description: 未找到需要更新的记录" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "未找到需要更新的记录",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "更新失败",
		})
		return
	}
}

func Delete(ctx *gin.Context) {
	process_product_id := ctx.PostForm("process_product_id")
	chaincode_name := "processcc"
	fnc := "delete"
	args := [][]byte{[]byte(process_product_id)}
	_, err := process_sdk.ChannelExecute(chaincode_name, fnc, args)
	if err!=nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "删除成功",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "删除失败",
	})
	return
}