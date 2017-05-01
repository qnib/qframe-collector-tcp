FROM qnib/uplain-golang

WORKDIR /usr/local/src/github.com/qnib/qframe-collector-tcp
COPY main.go ./main.go
COPY lib/ ./lib/
COPY vendor/vendor.json ./vendor/vendor.json
RUN govendor fetch +m \
 && govendor build

FROM qnib/uplain-init

COPY --from=0 /usr/local/src/github.com/qnib/qframe-collector-tcp/qframe-collector-tcp \
     /usr/local/bin/
ENV SKIP_ENTRYPOINTS=true
CMD ["qframe-collector-tcp"]
