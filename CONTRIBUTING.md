# Parvares

## Thank You :tada:

First off, thank you for reading this contributing document. It means a lot to us that
you're interested and are considering helping this project grow.

## Running Parvaeres on your local machine

It is very useful to run Parvaeres locally if you intend to contribute. In this way, you
can test its functionality and check your changes immediately.

You can use any kubernetes cluster with [ArgoCD](https://argoproj.github.io/argo-cd/)
installed to run Parvaeres server. To make things easier, you can use the script
`scripts/k3s-test-cluster.sh` to spin up a local [k3s](https://k3s.io/) cluster with
ArgoCD.

```shell
scripts/k3s-test-cluster.sh up
```

After this, you can deploy the Parvares Server using the manifests in the `deploy`
directory.

```shell
kubectl apply -f deploy/parvaeres-server.yaml
```

## Testing

You can see Parvares in action by deploying a simple application provided by the ArgoCD
project, but you can use any available kuberenetes manifests directory shared in a public
repository.

Before talking to the server you need to retrieve the its IP address:
```
export PARVARES_SERVER_IP=$(kubectl get services -n argocd parvaeres-server -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

Then, you'll able to send a request like this:
```
curl -H 'Content-Type: application/json' \
    -d '{"Email":"whatever@example.com","Repository":"https://github.com/argoproj/argocd-example-apps.git","Path":"guestbook"}'\
    http://$PARVARES_SERVER_IP:8080/v1/deployment -v
```

The response should be something similar:
```
* Expire in 0 ms for 6 (transfer 0x5560e788bf50)
*   Trying 172.25.0.3...
* TCP_NODELAY set
* Expire in 200 ms for 4 (transfer 0x5560e788bf50)
* Connected to 172.25.0.3 (172.25.0.3) port 8080 (#0)
> POST /v1/deployment HTTP/1.1
> Host: 172.25.0.3:8080
> User-Agent: curl/7.64.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 119
> 
* upload completely sent off: 119 out of 119 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Tue, 08 Sep 2020 13:38:14 GMT
< Content-Length: 99
< 
{"Message":"CREATED","Items":[{"UUID":"2c17e346-f6c2-4698-8864-57970b9380c1","Status":"PENDING"}]}
* Connection #0 to host 172.25.0.3 left intact
```

At this point the application is created but not yet deployed (`PENDING`). To "confirm"
the deployment you will need to visit the deployment page:

```
curl http://172.25.0.5:8080/v1/deployment/2c17e346-f6c2-4698-8864-57970b9380c1
{"Message":"FOUND","Items":[{"UUID":"2c17e346-f6c2-4698-8864-57970b9380c1","Status":"DEPLOYED"}]}
```

The deployment is now in `DEPLOYED` state, if everything went well. You should be able to
see the corresponding resources deployed in the namespace `2c17e346-f6c2-4698-8864-57970b9380c1`:

```
➜ kubectl get all -n 2c17e346-f6c2-4698-8864-57970b9380c1
NAME                                READY   STATUS    RESTARTS   AGE
pod/guestbook-ui-85c9c5f9cb-9gh76   1/1     Running   0          88s

NAME                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/guestbook-ui   ClusterIP   10.43.172.120   <none>        80/TCP    88s

NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/guestbook-ui   1/1     1            1           88s

NAME                                      DESIRED   CURRENT   READY   AGE
replicaset.apps/guestbook-ui-85c9c5f9cb   1         1         1       88s
```

**NOTE**: The UUID will be secret at some point and will not be exposed in the API call
response. It will be sent by email to the user. In this way, only the owner of the email
would be able to confirm and access the deployment. This security mechanism might prove
not sufficient, but it would be nice if we could keep it as simple as possible.
