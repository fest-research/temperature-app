FROM floreks/kubepi-base

COPY build/temperature-app-linux-arm-6 /usr/bin/temperature-app

ENTRYPOINT ["/usr/bin/temperature-app"]
