RemoteHost = 47.120.77.40
APPDIR=$(CURDIR)/app
BACKENDDIR=$(CURDIR)/backend
WEBDIR=$(CURDIR)/web
DEPLOYDIR=$(CURDIR)/deploy

AppBin=app_bin
BackendBin=backend_bin
webTar=web.tar


run_app:
	cd ${APPDIR}/cmd && go run main.go -c ${DEPLOYDIR}/app_debug_conf.yaml

run_backend:
	cd ${BACKENDDIR} && go run main.go -c ${DEPLOYDIR}/backend_debug_conf.yaml

run_web:
	cd ${WEBDIR} && npm run serve

deploy_app:
	cd ${APPDIR} && GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/${AppBin} ${APPDIR}/cmd/main.go
	sshpass -p 1qaz@WSX scp $(CURDIR)/${AppBin} root@${RemoteHost}:/root/caishuji

deploy_backend:
	cd ${BACKENDDIR} && GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/${BackendBin} ${BACKENDDIR}/main.go
	sshpass -p 1qaz@WSX scp $(CURDIR)/${BackendBin} root@${RemoteHost}:/root/caishuji

deploy_all:
	sshpass -p 1qaz@WSX  ssh -o StrictHostKeyChecking=no root@${RemoteHost}  "cd /root/caishuji && ./deploy2.sh"

# deploy_web:
# 	tar -cvf ${webTar} ./web
# 	sshpass -p 1qaz@WSX scp $(CURDIR)/${webTar} root@${RemoteHost}:/root/caishuji
# 	sshpass -p 1qaz@WSX  ssh -o StrictHostKeyChecking=no root@${RemoteHost}  "cd /root/caishuji && tar -xvf web.tar"
