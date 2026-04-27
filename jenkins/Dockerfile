FROM jenkins/jenkins:lts-jdk17

USER root

# Install Docker CLI
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
    $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && apt-get install -y docker-ce-cli && \
    rm -rf /var/lib/apt/lists/*

# Create docker group (standard GID 999) and add jenkins user
# We use || true in case the group already exists
RUN groupadd -g 999 docker || true && \
    usermod -aG docker jenkins

# Environment variable to disable setup wizard
ENV JAVA_OPTS="-Djenkins.install.runSetupWizard=false"

USER jenkins

# Install required and extra plugins
RUN jenkins-plugin-cli --plugins \
    workflow-aggregator \
    docker-workflow \
    pipeline-stage-view \
    git \
    configuration-as-code \
    sonar \
    blueocean \
    ansicolor \
    junit \
    cobertura \
    htmlpublisher \
    slack \
    warnings-ng
