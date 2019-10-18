IMAGE := ksync/ksync-example-frontend:latest

serve:
	npm run start

build:
	docker build -t $(IMAGE) .

push: build
	docker push $(IMAGE)
