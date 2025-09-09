FROM golang:1.25-alpine AS builder

WORKDIR /build
#Copying files from the source (our host machine) to the container destination - build ditrectory
#Copying ALL files in current directory to build directory 
#COPY <src> <dist>
COPY . .
RUN go mod download
#name of image: docker-trial
RUN go build -o ./docker-trial

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/docker-trial ./docker-trial
EXPOSE 8080
CMD [ "/app/docker-trial" ]