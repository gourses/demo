FROM scratch
COPY /target/bin/demo /
ENTRYPOINT ["/demo"]