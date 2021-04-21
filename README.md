# tryout-dynamic-cr-informer

An example to try out using Kubernetes dynamic informer to watch custom resource objects.

## Preparation

1. Ensure you have a Kubernets cluster around
1. Install the CRDs
    ```
    $ k apply -f manifests/foo-namespaced-crd.yaml
    $ k apply -f manifests/bar-clustered-crd.yaml
    ```

## Run the example

- go build -o main *.go
- ./main --kubeconfig=$HOME/.kube/config (keep this terminal opened)

Then in another terminal, try out changes on CR objects:

```
$ k apply -f manifests/foo1.yaml
$ k apply -f manifests/bar1.yaml
```

You should be able to see logs in the original terminal:

```
I0421 11:12:13.808740   29231 main.go:40] Started.
I0421 11:12:17.339746   29231 main.go:60] received "foo" add event!
I0421 11:12:24.255469   29231 main.go:60] received "bar" add event!
```
