package main

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"strings"
	"time"
)

type QualityCc struct {
}

func (p *QualityCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	quality_product_id := args[0]     //检测产品标识
	quality_date := args[1]        	  //检测日期
	product_quality := args[2]        //检测产品质量是否合格
	quality_job_id := args[3]         //检测审核人工号
	quality_job_name := args[4]       //检测审核人姓名
	txid :=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"quality_product_id":"%s","quality_date":"%s","product_quality":"%s","quality_job_id":"%s","quality_job_name":"%s","timestamp":"%s","txid":"%s"}`, quality_product_id,quality_date, product_quality,quality_job_id,quality_job_name,timestamp,txid)
	err = stub.PutState(quality_product_id, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *QualityCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args = stub.GetFunctionAndParameters()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	quality_product_id := args[0]     //检测产品标识
	quality_date := args[1]        	  //检测日期
	product_quality := args[2]        //检测产品质量是否合格
	quality_job_id := args[3]         //检测审核人工号
	quality_job_name := args[4]       //检测审核人姓名
	txid :=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(quality_product_id)
	if string(rsp) != "" {
		return shim.Error("该产品标识号已存在")
	}
	value := fmt.Sprintf(`{"quality_product_id":"%s","quality_date":"%s","product_quality":"%s","quality_job_id":"%s","quality_job_name":"%s","timestamp":"%s","txid":"%s"}`, quality_product_id,quality_date, product_quality,quality_job_id,quality_job_name,timestamp,txid)
	err = stub.PutState(quality_product_id, []byte(value))
	if err != nil {
		return shim.Error("新增产品失败")
	}
	return shim.Success([]byte("新增产品成功"))
}

func (p *QualityCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	quality_product_id := args[0]
	quality_product_id_byte, err := stub.GetState(quality_product_id)
	if err != nil {
		return shim.Error("没有产品信息")
	}
	return shim.Success(quality_product_id_byte)
}

func (p *QualityCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	quality_product_id := args[0]
	rsp,err:=stub.GetState(quality_product_id)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该产品编号不存在")
	}
	err = stub.DelState(quality_product_id)
	if err != nil {
		return shim.Error("没有待删除产品信息")
	}
	return shim.Success([]byte("删除质检产品成功"))
}

func (p *QualityCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args = stub.GetFunctionAndParameters()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	quality_product_id := args[0]     //检测产品标识
	quality_date := args[1]        	  //检测日期
	product_quality := args[2]        //检测产品质量是否合格
	quality_job_id := args[3]         //检测审核人工号
	quality_job_name := args[4]       //检测审核人姓名
	txid :=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(quality_product_id)
	if err!=nil ||result ==nil {
		return shim.Error("未找到需要更新的记录")
	}
	value := fmt.Sprintf(`{"quality_product_id":"%s","quality_date":"%s","product_quality":"%s","quality_job_id":"%s","quality_job_name":"%s","timestamp":"%s","txid":"%s"}`, quality_product_id,quality_date, product_quality,quality_job_id,quality_job_name,timestamp,txid)
	err = stub.PutState(quality_product_id, []byte(value))
	if err != nil {
		return shim.Error("更新质检记录失败")
	}
	return shim.Success([]byte("更新质检记录成功"))
}

func (p *QualityCc) queryQulityHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	quality_product_id :=args[0]
	// 获取历史变更数据，返回一个迭代数组,
	//the historic value and associated transaction id and timestamp are returned.
	iterator, err := stub.GetHistoryForKey(quality_product_id)
	if err != nil {
		return shim.Error("根据指定的标识号查询对应的历史变更数据失败")
	}
	defer iterator.Close()
	// 迭代处理
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		hisData, err := iterator.Next()//获取当前i的的迭代器数据
		if err != nil {
			return shim.Error("获取历史变更数据失败")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(";;")
		}
		buffer.WriteString(string(hisData.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *QualityCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	if f == "set" {
		return p.setValue(stub, args)
	}
	if f == "query" {
		return p.query(stub, args)
	}
	if f == "delete" {
		return p.delete(stub, args)
	}
	if f == "update" {
		return p.update(stub, args)
	}
	if f == "queryQulityHistoryByID" {
		return p.queryQulityHistoryByID(stub, args)
	}
	return shim.Error("函数名称错误，只能是set、query、delete、update或者queryQulityDetailByID")
}

func main() {
	err := shim.Start(new(QualityCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}
