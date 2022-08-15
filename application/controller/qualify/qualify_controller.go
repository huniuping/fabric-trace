package quality

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	quality_sdk "raft-fabric-project/application/sdk/quality"
	"strings"
)

func Query(ctx *gin.Context) {
	quality_product_id := ctx.Query("quality_product_id")
	chaincode_name := "qualitycc"
	fnc := "query"
	args := [][]byte{[]byte(quality_product_id)}
	rsp, err := quality_sdk.ChannelQuery(chaincode_name, fnc, args)
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

func QueryQulityDetailByID(ctx *gin.Context) {
	quality_product_id := ctx.Query("quality_product_id")
	chaincode_name := "qualitycc"
	fnc := "queryQulifyDetailByID"
	args := [][]byte{[]byte(quality_product_id)}
	rsp, err := quality_sdk.ChannelQuery(chaincode_name, fnc, args)
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
	qualitys:=[]map[string]interface{}{}

	for i := 0; i < len(payloads); i++ {
		if payloads[i] ==""{
			continue
		}
		temp_map := make(map[string]interface{})//一行所有数据存在这一个map中
		json.Unmarshal([]byte(payloads[i]), &temp_map)
		qualitys=append(qualitys,temp_map)
	}
	map_data["qualityHistory"]=qualitys
	fmt.Println(map_data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"meg":  "查询成功",
		"data": map_data,
	})
	return
}

func Delete(ctx *gin.Context) {
	quality_product_id := ctx.PostForm("quality_product_id")
	chaincode_name := "qualitycc"
	fnc := "delete"
	args := [][]byte{[]byte(quality_product_id)}
	_, err := quality_sdk.ChannelExecute(chaincode_name, fnc, args)
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

func Set(ctx *gin.Context) {
	quality_product_id := ctx.PostForm("quality_product_id")
	quality_date := ctx.PostForm("quality_date")
	product_quality := ctx.PostForm("product_quality")
	quality_job_id := ctx.PostForm("quality_job_id")
	quality_job_name := ctx.PostForm("quality_job_name")
	chaincode_name := "qualitycc"
	fnc := "set"
	args := [][]byte{[]byte(quality_product_id),
		[]byte(quality_date),
		[]byte(product_quality),
		[]byte(quality_job_id),
		[]byte(quality_job_name)}
	_, err := quality_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经添加成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:13051]: Chaincode status Code: (500) UNKNOWN. Description: 该产品标识号已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该产品已存在",
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
	quality_product_id := ctx.PostForm("quality_product_id")
	quality_date := ctx.PostForm("quality_date")
	product_quality := ctx.PostForm("product_quality")
	quality_job_id := ctx.PostForm("quality_job_id")
	quality_job_name := ctx.PostForm("quality_job_name")

	chaincode_name := "qualitycc"
	fnc := "update"
	args := [][]byte{[]byte(quality_product_id),
		[]byte(quality_date),
		[]byte(product_quality),
		[]byte(quality_job_id),
		[]byte(quality_job_name)}
	_, err := quality_sdk.ChannelExecute(chaincode_name, fnc, args)
	fmt.Println("====================")
	fmt.Println(err)
	if error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已经更新成功",
		})
		return
	} else if error.Error(err) == "Transaction processing for endorser [localhost:13051]: Chaincode status Code: (500) UNKNOWN. Description: 该产品标识号已存在" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "该产品已存在",
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