FROM golang:1.25.3-trixie

WORKDIR "/app"
RUN go install golang.org/x/tools/cmd/goimports@latest

CMD ["/bin/bash"]
