FROM jenkins/jenkins:lts-jdk17

USER root

# Install Docker CLI (DooD approach with host socket)
RUN apt-get update && apt-get install -y \
    docker.io \
    git && \
    rm -rf /var/lib/apt/lists/*

# Ensure jenkins user is in docker group to access /var/run/docker.sock
RUN getent group docker || groupadd docker && usermod -aG docker jenkins

USER jenkins

# Pre-install challenge-required plugins
RUN jenkins-plugin-cli --plugins \
    workflow-aggregator \
    git \
    docker-workflow \
    pipeline-stage-view \
    configuration-as-code \
    job-dsl \
    sonar

USER jenkins
