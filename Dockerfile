FROM gcr.io/distroless/static:nonroot

COPY _output/soubise /usr/local/bin/

ENTRYPOINT ["soubise"]
