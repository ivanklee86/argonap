# See here for image contents: https://github.com/microsoft/vscode-dev-containers/blob/main/containers/go/.devcontainer/Dockerfile
ARG VARIANT="1"
FROM mcr.microsoft.com/devcontainers/go:1-${VARIANT}

# Install dependencies.
# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

USER vscode

# Install tools
# - Tools
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
COPY --from=quay.io/argoproj/argocd /usr/local/bin/argocd /usr/local/bin

# - Golang
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/jstemmer/go-junit-report/v2@latest
RUN go install github.com/hairyhenderson/gomplate/v4/cmd/gomplate@latest

# - Kubernetes
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

COPY --from=bitnami/kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin
COPY --from=alpine/helm /usr/bin/helm /usr/local/bin
COPY --from=derailed/k9s /bin/k9s /usr/local/bin

