package contract

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"net/http"
	contract_sdk "raft-fabric-project/application/sdk/contract"
	"strings"
)

func Query(ctx *gin.Context) {
	contractId := ctx.Query("contractId")
	chaincode_name := "contractcc"
	fnc := "query"
	args := [][]byte{[]byte(contractId)}
	rsp, err := contract_sdk.ChannelQuery(chaincode_name, fnc, args)
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

func QueryByProductName(ctx *gin.Context) {
	productName := ctx.Query("productName")
	chaincode_name := "contractcc"
	fnc := "queryByProductName"
	args := [][]byte{[]byte(productName)}
	rsp, err := contract_sdk.ChannelQuery(chaincode_name, fnc, args)
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

func QueryContractDetailByID(ctx *gin.Context) {
	contractId := ctx.Query("contractId")
	chaincode_name := "contractcc"
	fnc := "queryContractDetailByID"
	args := [][]byte{[]byte(contractId)}
	rsp, err := contract_sdk.ChannelQuery(chaincode_name, fnc, args)
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
		if payloads[i] ==""{
			continue
		}
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		contracts=append(contracts,temp_map)
	}
	map_data["contractsHistory"]=contracts
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func Delete(ctx *gin.Context) {
	contractId := ctx.PostForm("contractId")
	chaincode_name := "contractcc"
	fnc := "delete"
	args := [][]byte{[]byte(contractId)}
	_, err := contract_sdk.ChannelExecute(chaincode_name, fnc, args)
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

func Querycontractfile(ctx *gin.Context){
	contract_hash := ctx.Query("contract_hash")
	contract_file:=CatIPFS(contract_hash)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": contract_file,
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
	contractId := ctx.PostForm("contractId")
	contractName := ctx.PostForm("contractName")
	buyer := ctx.PostForm("buyer")
	seller := ctx.PostForm("seller")
	productName := ctx.PostForm("productName")
	productAmount := ctx.PostForm("productAmount")
	fileheader, err := ctx.FormFile("contract_hash")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	contract_hash := UploadIPFS(string(file))
	contract_file := contract_hash
	chaincode_name := "contractcc"
	fnc := "set"
	args := [][]byte{[]byte(contractId),
		[]byte(contractName),
		[]byte(buyer),
		[]byte(seller),
		[]byte(productName),
		[]byte(productAmount),
		[]byte(contract_file)}
	_, err = contract_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经添加成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:17051]: Chaincode status Code: (500) UNKNOWN. Description: 该订单编号-期数已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该合同单号已存在",
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
	contractId := ctx.PostForm("contractId")
	contractName := ctx.PostForm("contractName")
	buyer := ctx.PostForm("buyer")
	seller := ctx.PostForm("seller")
	productName := ctx.PostForm("productName")
	productAmount := ctx.PostForm("productAmount")
	fileheader, err := ctx.FormFile("contract_hash")
	file, err := json.Marshal(&fileheader)
	if err != nil {
		fmt.Println("序列化err=", err)
	}
	contract_hash := UploadIPFS(string(file))
	contract_file := contract_hash
	chaincode_name := "contractcc"
	fnc := "update"
	args := [][]byte{[]byte(contractId),
		[]byte(contractName),
		[]byte(buyer),
		[]byte(seller),
		[]byte(productName),
		[]byte(productAmount),
		[]byte(contract_file)}
	_, err = contract_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经更新成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:17051]: Chaincode status Code: (500) UNKNOWN. Description: 未找到需要更新的记录" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该合同单号已存在",
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
