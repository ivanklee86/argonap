# Shared Go base (ivanklee86/devcontainers). It already provides golangci-lint,
# go-junit-report, gomplate, goreleaser, prek, Task, gh and Claude Code.
FROM ghcr.io/ivanklee86/devcontainer/go:1.27

# Kubernetes / Argo CD tooling for the local k3d test stack. The versions are kept
# in step with tangle's .devcontainer/Dockerfile.
ARG K3D_VERSION=v5.9.0 # github-releases/k3d-io/k3d
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=${K3D_VERSION} bash

COPY --from=rancher/kubectl:v1.36.2 /bin/kubectl /usr/local/bin/kubectl
COPY --from=alpine/helm:4.3.2 /usr/bin/helm /usr/local/bin
COPY --from=derailed/k9s:v0.50.18 /bin/k9s /usr/local/bin
COPY --from=quay.io/argoproj/argocd:v3.6.3 /usr/local/bin/argocd /usr/local/bin
