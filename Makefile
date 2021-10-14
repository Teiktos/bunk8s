PROJECT = bunk8s

build-coordinator:

	cd coordinator/src/main && docker build . -t $(PROJECT)-coordinator-image:latest

build-launcher:

	cd launcher/src/main && docker build . -t $(PROJECT)-launcher-image:latest

build: build-coordinator build-launcher 

start-coordinator: build

	cd coordinator/helm && helm install $(PROJECT)-coordinator $(PROJECT)-coordinator

stop-coordinator:

	helm uninstall $(PROJECT)-coordinator

start-launcher: build-launcher

	docker run $(PROJECT)-launcher-image

stop-launcher:

	docker stop 

restart-coordinator: stop-coordinator start-coordinator


		
