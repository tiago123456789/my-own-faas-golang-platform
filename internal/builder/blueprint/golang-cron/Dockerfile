ARG VERSION_TAG

FROM ${VERSION_TAG}-alpine as builder

ENV PATH_CODE=../lambda

ARG MODULE_PATH

WORKDIR app/app

RUN mkdir ../lambda

COPY main.go ./

COPY code.zip ../lambda

COPY go.mod go.sum ./

RUN apk add --no-cache unzip

RUN cd ../lambda && unzip code.zip

RUN cd ../lambda && go mod tidy

RUN sed -i 's|MODULE_PATH|'$MODULE_PATH'|' go.mod

RUN sed -i 's|MODULE_NAME|'$MODULE_PATH'|' main.go

RUN go mod tidy

RUN go build -ldflags "-s -w" -o application ./

FROM alpine:latest  

ENV PATH_CODE=../lambda

COPY --from=builder /go/app/app/application .

RUN mkdir lambda 

COPY --from=builder /go/app/lambda/config.yml ./lambda

CMD ["./application"]
