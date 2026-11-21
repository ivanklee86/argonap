{
  name: "argonap",
  // Argo CD, served through k3d's load balancer (see k3d-config.yaml).
  forwardPorts: [8080],
  postCreateCommand+: {
    "argonap-post-install": "bash ./.devcontainer/post_install.sh",
  },
}
