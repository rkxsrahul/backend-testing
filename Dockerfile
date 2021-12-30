FROM golang:1.13 as golang
MAINTAINER Gursimran Singh <singhgursimran@me.com>

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin
ARG BUILD_ID
ENV BUILD_IMAGE=$BUILD_ID
ENV GO111MODULE=off

# build directories
ADD . /go/src/git.xenonstack.com/util/continuous-security-backend
WORKDIR /go/src/git.xenonstack.com/util/continuous-security-backend

RUN go install git.xenonstack.com/util/continuous-security-backend
EXPOSE 8000

# new build stage started and copy artifacts from previous stage
#FROM python:3.7-stretch
FROM nikolaik/python-nodejs:python3.6-nodejs13-stretch
ARG BUILD_ID
ENV BUILD_IMAGE=$BUILD_ID

# install dependency packages
RUN apt update; apt-get install -y nmap whois dnsutils ldnsutils ldnsutils bsdmainutils netcat openssl ca-certificates
#RUN npm install wappalyzer -g --unsafe-perm

# clone git repository
RUN git clone https://github.com/meliot/shcheck.git $HOME/projects/security/tools/shcheck; git clone https://github.com/nixcraft/domain-check-2.git $HOME/projects/security/tools/domain-check ; wget https://github.com/drwetter/testssl.sh/archive/refs/tags/3.0.5.tar.gz ; tar xvf 3.0.5.tar.gz -C $HOME/projects/security/tools/ ; wget https://raw.githubusercontent.com/nixcraft/domain-check-2/master/domain-check-2.sh ; mv domain-check-2.sh /usr/local/bin/domain-check-2.sh ; chmod +x /usr/local/bin/domain-check-2.sh 
RUN cp -r $HOME/projects/security/tools/testssl.sh-3.0.5  $HOME/projects/security/tools/testssl.sh

#create folder to copy code from previous stage and executable from previous stage
RUN mkdir -p /go/src/git.xenonstack.com/util/continuous-security-backend; mkdir -p /go/bin

#set working directory
WORKDIR /go/src/git.xenonstack.com/util/continuous-security-backend

#copy code from previous stage
COPY --from=golang  /go/src/git.xenonstack.com/util/continuous-security-backend /go/src/git.xenonstack.com/util/continuous-security-backend

#copy executable file from previous stage
COPY --from=golang /go/bin/continuous-security-backend  /go/bin/

#clonning scripts
RUN export GIT_SSL_NO_VERIFY=1 && git clone https://gitlab-ci-token:LisfzisY1Ly2oxmWGiBJ@git.xenonstack.com/devops/web-security.git -b develop --single-branch tools/
# RUN git clone https://gitlab-ci-token:LisfzisY1Ly2oxmWGiBJ@git.xenonstack.com/devops/web-security.git -b develop --single-branch tools/
RUN cp -r /go/src/git.xenonstack.com/util/continuous-security-backend/tools/* $HOME/projects/security/tools/
RUN mv /go/src/git.xenonstack.com/util/continuous-security-backend/tools /go/src/git.xenonstack.com/util/continuous-security-backend/scripts

#install domain-check2 tool
#RUN wget https://raw.githubusercontent.com/nixcraft/domain-check-2/master/domain-check-2.sh ; cp -vf domain-check-2.sh /usr/local/bin/domain-check-2.sh ; chmod +x /usr/local/bin/domain-check-2.sh
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.20.2


#liston on port
EXPOSE 8000
