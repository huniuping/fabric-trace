# fabric-trance
The project currently provides the building files, SDK files, business interface files and blockchain network startup scripts of the fabric blockchain network.

If you want to clone the project in reverse and start it locally,you should make the cloned project path as same as this project ,or you need to modify the certificate path and chaincode path in some yaml configuration files in the fabric network to the file path you generated.

In addition, the network startup script in the project is a summary script, which is cloned locally for the first time. It is recommended to open the script file and execute the script command step by step to view the error location.

Here are some tips for authors when creating projects：
1.Environment for this project: Ubuntu 18.04, fabric 2.3, docker 20.10.7, docker compose 1.17.1 and go1.17.8 Linux / AMD64
2.You need to add DNS mapping for container on host, because this project is built on a single machine
3.In the chaincode folder of fabric2.3, in addition to the chaincode files, the vendor directory and go.mod, go The author suggests that “go mod” be used to manage the dependency packages of chain code during chaincode development. The vendor can obtain them through the command “go mod vendor”

Specially,the SASproject file is the code that implement ordered aggregate signature ,it's coded by java,and it used jpbc library to support the signature.
