FROM golang:1.18

RUN apt update && apt install -y

# Language autocomplete features.
RUN go install golang.org/x/tools/gopls@latest
RUN go get golang.org/x/tools/cmd/goimports

CMD /usr/bin/bash
