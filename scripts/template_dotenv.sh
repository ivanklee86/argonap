#!/bin/sh

argocd account generate-token --account automation | gomplate -f .env.tmpl -d token=stdin: > .env
