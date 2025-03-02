FROM oryd/kratos:v1.3.1-distroless

COPY kratos.yaml /etc/config/kratos/kratos.yaml
COPY identity.schema.json /etc/config/kratos/identity.schema.json
COPY github.jsonnet /etc/config/kratos/github.jsonnet
COPY webhook.jsonnet /etc/config/kratos/webhook.jsonnet


ENTRYPOINT [ "kratos" ]

CMD [ "serve",  "-c", "/etc/config/kratos/kratos.yaml", "--watch-courier" ]
