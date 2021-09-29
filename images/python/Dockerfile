FROM amazonlinux

ENV LAMBDA_RUNTIME=python
ENV LAMBDA_BOOTSTRAP_BINARY=/lambda/bootstrap.py
ENV BOOTSTRAP_LOG_FILE=/lambda/bootstrap.log

WORKDIR /lambda/lib/
COPY bin/runtime /lambda/lib/runtime
COPY images/python/entrypoint.sh /lambda/lib/entrypoint.sh
COPY images/python/bootstrap.py /lambda/bootstrap.py

WORKDIR /lambda/code/
COPY images/python/__init__.py /lambda/code/__init__.py

RUN yum install python3 unzip -y pip3 install boto3 && chmod 777 /lambda/bootstrap.py && chmod 777 /lambda/lib/entrypoint.sh
ENTRYPOINT [ "/lambda/lib/entrypoint.sh" ]