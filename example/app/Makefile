IMAGE := ksync/ksync-example:latest

serve:
	python -m flask run --port=8080

build:
	docker build -t $(IMAGE) .

push: build
	docker push $(IMAGE)
