FROM scratch

COPY fxtx /bin/

ENTRYPOINT ["/bin/fxtx"]