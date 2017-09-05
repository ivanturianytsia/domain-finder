FROM alpine

RUN apk --update upgrade && \
    apk --update add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ENV APP_NAME="domainfinder"
ENV PARENT_DIR="/go/src/github.com/ivanturianytsia/"${APPNAME}
ENV PORT="8000"


RUN mkdir -p ${PARENT_DIR}
RUN mkdir -p ${PARENT_DIR}/app

WORKDIR ${PARENT_DIR}

ADD ./bin/${APP_NAME}_alpine ${PARENT_DIR}/${APP_NAME}
ADD ./app ${PARENT_DIR}/app

CMD ${PARENT_DIR}"/"${APP_NAME}
