FROM rabbitmq:4.2

RUN apt-get update && apt-get install -y gettext-base

RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_prometheus

COPY ./init.sh /usr/local/bin/init.sh
RUN chmod +x /usr/local/bin/init.sh

CMD ["/usr/local/bin/init.sh"]
