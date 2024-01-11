FROM golang:1.21

RUN apt update && apt install -y

# Language autocomplete features.
RUN go install golang.org/x/tools/gopls@latest
RUN go install golang.org/x/tools/cmd/goimports@latest

CMD /usr/bin/bash
