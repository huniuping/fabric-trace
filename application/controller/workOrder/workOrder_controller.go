package workOrder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"net/http"
	manufacture_sdk "raft-fabric-project/application/sdk/manufacture"
	"strings"
)

func Query(ctx *gin.Context) {
	work_order_id := ctx.Query("work_order_id")
	chaincode_name := "gongdancc"
	fnc := "query"
	args := [][]byte{[]byte(work_order_id)}
	rsp, err := manufacture_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
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

func QueryByContractId(ctx *gin.Context) {
	contractId := ctx.Query("contractId")
	chaincode_name := "gongdancc"
	fnc := "queryByContractId"
	args := [][]byte{[]byte(contractId)}
	rsp, err := manufacture_sdk.ChannelQuery(chaincode_name, fnc, args)
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
	contracts:=[]map[string]interface{}{}
	for i := 0; i < len(payloads); i++ {
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		contracts=append(contracts,temp_map)
	}
	map_data["contracts"]=contracts
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func QueryByDrawingID(ctx *gin.Context) {
	drawing_id := ctx.Query("drawing_id")
	chaincode_name := "gongdancc"
	fnc := "queryByDrawingID"
	args := [][]byte{[]byte(drawing_id)}
	rsp, err := manufacture_sdk.ChannelQuery(chaincode_name, fnc, args)
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
	drawings:=[]map[string]interface{}{}

	for i := 0; i < len(payloads); i++ {
		if payloads[i] ==""{
			continue
		}
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		drawings=append(drawings,temp_map)
	}
	map_data["drawings"]=drawings
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func Queryworkorderfile(ctx *gin.Context){
	workorder_hash := ctx.Query("workorder_hash")
	workorder_file:=CatIPFS(workorder_hash)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": workorder_file,
	})
	return
}

func UploadIPFS(str string) string {
	sh = shell.NewShell("localhost:5001")

	hash, err := sh.Add(bytes.NewBufferString(str))
	if err != nil {
		fmt.Println("上传ipfs时错误：", err)
	}
	return hash
}

func CatIPFS(hash string) string {
	sh = shell.NewShell("localhost:5001")
	read, err := sh.Cat(hash)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(read)
	return string(body)
}

func Set(ctx *gin.Context) {
	work_order_id := ctx.PostForm("work_order_id")
	duration := ctx.PostForm("duration")
	drawing_id := ctx.PostForm("drawing_id")
	contract_id := ctx.PostForm("contract_id")
	fileheader, err := ctx.FormFile("task")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	task_hash := UploadIPFS(string(file))
	task := task_hash
	chaincode_name := "gongdancc"
	fnc := "set"
	args := [][]byte{[]byte(work_order_id),
		[]byte(task),
		[]byte(duration),
		[]byte(drawing_id),
		[]byte(contract_id)}
	_, err = manufacture_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经添加成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:19051]: Chaincode status Code: (500) UNKNOWN. Description: 该工单标识号已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该工单编号已存在",
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
	work_order_id := ctx.PostForm("work_order_id")
	duration := ctx.PostForm("duration")
	drawing_id := ctx.PostForm("drawing_id")
	contract_id := ctx.PostForm("contract_id")
	fileheader, err := ctx.FormFile("task")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	task_hash := UploadIPFS(string(file))
	task := task_hash
	chaincode_name := "gongdancc"
	fnc := "update"
	args := [][]byte{[]byte(work_order_id),
		[]byte(task),
		[]byte(duration),
		[]byte(drawing_id),
		[]byte(contract_id)}
	_, err = manufacture_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经更新成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:7051]: Chaincode status Code: (500) UNKNOWN. Description: 未找到需要更新的图纸" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "未找到需要更新的图纸",
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
	work_order_id := ctx.PostForm("work_order_id")
	chaincode_name := "gongdancc"
	fnc := "delete"
	args := [][]byte{[]byte(work_order_id)}
	_, err := manufacture_sdk.ChannelExecute(chaincode_name, fnc, args)
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



