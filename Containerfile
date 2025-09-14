FROM golang AS build-stage

# Copy source code
ADD box /box
ADD boxer /boxer

# Compile uvbox CLI
WORKDIR /boxer
RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /uvbox

# Release in small distroless container
FROM gcr.io/distroless/static-debian12 AS build-release-stage
COPY --from=build-stage /uvbox /uvbox
USER nonroot:nonroot
ENTRYPOINT ["/uvbox"]
