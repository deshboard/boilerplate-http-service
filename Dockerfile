FROM scratch

ARG BUILD_DIR
ARG BINARY_NAME

COPY $BUILD_DIR/$BINARY_NAME /service

EXPOSE 80 10000
CMD ["/service", "-http.addr", ":80", "-debug.addr", ":10000"]
HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost:10000/healthz || exit 1
