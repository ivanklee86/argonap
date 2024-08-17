#!/bin/sh

argocd account generate-token --account automation --grpc-web | gomplate -f .env.tmpl -d token=stdin: > .env
