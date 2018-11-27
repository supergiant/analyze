# Analyze 
![sg_analyze_light](https://user-images.githubusercontent.com/2936828/48772107-0b305300-eccc-11e8-8c72-4bcbd737226b.png)

[![Coverage Status](https://coveralls.io/repos/github/supergiant/analyze/badge.svg?branch=master)](https://coveralls.io/github/supergiant/analyze?branch=master)
[![Build Status](https://travis-ci.org/supergiant/analyze.svg?branch=master)](https://travis-ci.org/supergiant/analyze)
[![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/supergiant/analyze)](https://goreportcard.com/report/github.com/supergiant/analyze)




#deployment using HELM

0. helm chart can be found on ./helm/analyze folder  
1. ```helm install --debug ./helm/analyze/ --set cloudProviderType=aws --set aws.region=us-wild-wild-west--1 --set aws.accessKeyId=xxx --set aws.secretAccessKey=xxxx ``` if you need to configure ingress use flag ```--set ingress.enabled=true```
2. (optional) if there is no ingress controller nginx is installed on cluster install it using command ```helm install stable/nginx-ingress```
3. (optional) if there is RBAC on cluster we need to enable view for service account ```kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default```  

# dev info
This repository using go 1.11 modules instead of vendoring.
In order to enable it export GO111MODULE=on  

Service configuration can be stored using JSON, TOML, YAML, HCL(v1), and Java properties config file
Also can be configured using environment variables:

*AZ_LOGGING_LEVEL `debug`  
*AZ_LOGGING_FORMATTER `TXT`  


