ETCD_VERSION=3.5.16
REGISTRY=us-east1-docker.pkg.dev/tks-gcr-pub
IMAGE_PREFIX=etcd-operator
SHORT_NAME=etcd

#build-etcd-image:
#	docker build  -f hack/build/custom-etcd-image/Dockerfile -t ${REGISTRY}/${IMAGE_PREFIX}/${SHORT_NAME}:${ETCD_VERSION} --build-arg VERSION=${ETCD_VERSION} .


build-etcd-image:
	docker build -f hack/build/custom-etcd-image/Dockerfile \
        -t ${REGISTRY}/${IMAGE_PREFIX}/${SHORT_NAME}:v${ETCD_VERSION}  \
        --build-arg VERSION=${ETCD_VERSION} \
        .
push-etcd-image:
	docker push ${REGISTRY}/${IMAGE_PREFIX}/${SHORT_NAME}:v${ETCD_VERSION}

build-and-push: build-etcd-image push-etcd-image
