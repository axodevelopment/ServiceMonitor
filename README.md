# ServiceMonitor

The following will describe how to install a servicemonitor and have it point to your app.  The repo contains a basic go app and a basic service monitor implementation.

## How to enable user workload monitoring

To enable user defined cluster monitoring you need to create a config map with the value:
An example is located  /extern_resource/...

```bash
enableUserWorkload: true
```

You would need to deploy that config or update one there.
```bash
oc edit configmap cluster-monitoring-config -n openshift-monitoring 
```

Make sure you see the user worloads running
```bash
oc  get pod -n openshift-user-workload-monitoring
```

## How to deploy the service monitor

I put in a sample go app that creates a metric based upon a timer.  This just toggles a value between 0 and 1 but it gives a start.  The deployment files are available in kustomization/...  There is also a Dockerfile if you want.

For now you can view kustomziation/servicemonitor.yaml on what a file should look like.

```bash
oc apply -f servicemonitor.yaml
```

A quick review of ServiceMonitor

```bash
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: prometheus-example-monitor
  namespace: servicemonitor-a
spec:
  endpoints:
  - interval: 30s
    port: web           <- points to svc port name
    scheme: http
  selector:
    matchLabels:
      app: prometheus-example-app <- points to svc label

Maps to Service

apiVersion: v1
kind: Service
metadata:
  name: prometheus-example-app
  namespace: servicemonitor-a
  labels:
    app: prometheus-example-app <- this label
spec:
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
      name: web        <- this port name
  selector:
    app: prometheus-example-app
  type: ClusterIP
```

# How to Debug

Here are some things i've seen be problems with initial setups of a ServiceMonitor...

## labels

As mentioned above check that the servicemonitor points to the service and not the deployment/pod.

```bash
oc get service prometheus-example-app -n servicemonitor-a --show-labels
oc get service prometheus-example-app -n servicemonitor-a -o yaml
```

## targets

In OCP, you should see a Target under Admin -> Observe -> Targets.

You can also query the targets in prom by querying api/v1/targets.  You should see a 0 in droppedCounters but more importantly you should see your service being queried.  If it is in a dropped counter it is probably not able to find the service.

```bash
PROM_POD=$(oc get pods -n openshift-user-workload-monitoring -l app.kubernetes.io/name=prometheus -o name | head -1)
oc port-forward -n openshift-user-workload-monitoring $PROM_POD 9090:9090

curl localhost:9090/api/v1/targets
```

## thanos

Querying thanos can also help you identify what may be happening.

```bash
THANOS_QUERIER_HOST=$(oc -n openshift-monitoring get route thanos-querier -o jsonpath='{.spec.host}')

curl -k -H "Authorization: Bearer $(oc whoami -t)" "https://$THANOS_QUERIER_HOST/api/v1/query?query=my_app_up"
curl -k -H "Authorization: Bearer $(oc whoami -t)" "https://$THANOS_QUERIER_HOST/api/v1/query?query=my_app_requests_total"
```


## confirm logs
Checking the logs...  your pod names will be different then mine.

```bash
oc logs -n openshift-user-workload-monitoring prometheus-operator-7c7b77cbff-7v7km -c prometheus-operator
```

## follow up

First the deployment app should be making metrics available at /metrics and be available ideally done by a prom client.

2.) The ServiceMonitor uses label selectors which point to a service.

3.)you should be able to query the /metrics and get data back from the app.

4.) in the user workload pods I should see something like discovery manager scrape referencing the service monitor I deployed.

5.) Openshift metrics should show target as up and running

6.) if I portforward the openshift promethus user workload pods to port 9090:9090 I can query that prom endpoint /api/v1/targets and I should see my endpoint in there.  If it is under the droppedTargetCounts I know that while the ServiceMonitor is configured it is unable to reach that endpoint or getting an error along the way


Missing service labels
Mismatched label selectors
Wrong namespace configuration
Port name mismatches
Metrics endpoint accessibility