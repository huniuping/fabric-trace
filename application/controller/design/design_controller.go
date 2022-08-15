package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	shell "github.com/ipfs/go-ipfs-api"
	design_sdk "raft-fabric-project/application/sdk/design"
	"strings"
)

func Query(ctx *gin.Context) {
	drawing_id := ctx.Query("drawing_id")
	chaincode_name := "drawingcc"
	fnc := "query"
	args := [][]byte{[]byte(drawing_id)}
	rsp, err := design_sdk.ChannelQuery(chaincode_name, fnc, args)
	if err != nil||string(rsp.Payload)=="" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	fmt.Println("palyload=============")
	fmt.Println(rsp.Payload)
	fmt.Println("stringpalyload=============")
	xx:=string(rsp.Payload)
	fmt.Println("bytepalyload=============")
	fmt.Println([]byte(string(rsp.Payload)))

	fmt.Println(xx)


	fmt.Println("map=============")

	map_data := make(map[string]interface{})

	//map_data1:=map[string]interface{}{}
	//bcpxh:=[]map[string]interface{}{}
	//bcpxh = append(
	//	bcpxh,
	//	map_data,
	//	map_data,
	//	map_data,
	//	map_data,
	//	)
	//map_data1["bcpxh"] = bcpxh

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
	chaincode_name := "drawingcc"
	fnc := "queryByContractId"
	args := [][]byte{[]byte(contractId)}
	rsp, err := design_sdk.ChannelQuery(chaincode_name, fnc, args)
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
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		drawings=append(drawings,temp_map)
	}
	map_data["drawings"]=drawings
	fmt.Println("111111111111111111111111111")
	fmt.Println(map_data)
	fmt.Println("111111111111111111111111111")
	fmt.Println("222222222222222222222222222")
	//fmt.Println(xx[1])
	fmt.Println("22222222222222222222222222")

	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,

	})
	return
}

func QueryDrawingDetailByID(ctx *gin.Context) {
	drawing_id := ctx.Query("drawing_id")
	chaincode_name := "drawingcc"
	fnc := "queryDrawingDetailByID"
	args := [][]byte{[]byte(drawing_id)}
	rsp, err := design_sdk.ChannelQuery(chaincode_name, fnc, args)
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
	map_data["drawingsHistory"]=drawings
	fmt.Println("111111111111111111111111111")
	fmt.Println(map_data)
	fmt.Println("111111111111111111111111111")
	fmt.Println("222222222222222222222222222")
	//fmt.Println(xx[1])
	fmt.Println("22222222222222222222222222")

	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,

	})
	return
}

func Querydrawingfile(ctx *gin.Context){
	drawing_hash := ctx.Query("drawing_hash")
	workorder_file:=CatIPFS(drawing_hash)
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
	drawing_id := ctx.PostForm("drawing_id")
	drawing_name := ctx.PostForm("drawing_name")
	//drawing_file := ctx.PostForm("producer_id")
	contractId := ctx.PostForm("contractId")
	technology := ctx.PostForm("technology")
	fileheader, err := ctx.FormFile("drawing_file")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	drawing_hash := UploadIPFS(string(file))
	drawing_file := drawing_hash
	chaincode_name := "drawingcc"
	fnc := "set"
	args := [][]byte{[]byte(drawing_id),
		[]byte(drawing_name),
		[]byte(drawing_file),
		[]byte(contractId),
		[]byte(technology)}
	_, err = design_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经添加成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:7051]: Chaincode status Code: (500) UNKNOWN. Description: 该图纸标识号已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该图纸编号已存在",
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
	drawing_id := ctx.PostForm("drawing_id")
	drawing_name := ctx.PostForm("drawing_name")
	contractId := ctx.PostForm("contractId")
	technology := ctx.PostForm("technology")
	fileheader, err := ctx.FormFile("drawing_file")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	drawing_hash := UploadIPFS(string(file))
	drawing_file := drawing_hash
	chaincode_name := "drawingcc"
	fnc := "update"
	args := [][]byte{[]byte(drawing_id),
		[]byte(drawing_name),
		[]byte(drawing_file),
		[]byte(contractId),
		[]byte(technology)}
	_, err = design_sdk.ChannelExecute(chaincode_name, fnc, args)
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
	drawing_id := ctx.PostForm("drawing_id")
	chaincode_name := "drawingcc"
	fnc := "delete"
	args := [][]byte{[]byte(drawing_id)}
	_, err := design_sdk.ChannelExecute(chaincode_name, fnc, args)
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



