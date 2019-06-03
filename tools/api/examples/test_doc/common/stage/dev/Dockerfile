FROM centos:latest
RUN mkdir -p /data/apps/test_doc
WORKDIR /data/app/test_doc
COPY ./deploy_dev_test_doc_app.zip ./
RUN unzip -o ./deploy_dev_test_doc_app.zip
RUN ls
CMD ./main