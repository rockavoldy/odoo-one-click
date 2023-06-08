FROM --platform=amd64 golang:1.20-buster as build

WORKDIR /app
COPY . ./
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod tidy && go vet . && go build -ldflags="-s -w" -o odoo-one-click .

FROM --platform=amd64 ubuntu:jammy
COPY --from=build /app/odoo-one-click /odoo-one-click

WORKDIR /app
CMD ["bash"]