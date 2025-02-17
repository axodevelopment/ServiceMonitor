# ServiceMonitorz

 oc -n openshift-monitoring edit configmap cluster-monitoring-config

 oc -n openshift-user-workload-monitoring get pod


 oc apply -f servicemonitor.yaml

### Debug

# labels
oc get service prometheus-example-app -n servicemonitor-a --show-labels
oc get service prometheus-example-app -n servicemonitor-a -o yaml

# targets
PROM_POD=$(oc get pods -n openshift-user-workload-monitoring -l app.kubernetes.io/name=prometheus -o name | head -1)
oc port-forward -n openshift-user-workload-monitoring $PROM_POD 9090:9090

curl localhost:9090/api/v1/targets

# thanos
THANOS_QUERIER_HOST=$(oc -n openshift-monitoring get route thanos-querier -o jsonpath='{.spec.host}')

curl -k -H "Authorization: Bearer $(oc whoami -t)" "https://$THANOS_QUERIER_HOST/api/v1/query?query=my_app_up"
curl -k -H "Authorization: Bearer $(oc whoami -t)" "https://$THANOS_QUERIER_HOST/api/v1/query?query=my_app_requests_total"

# confirm logs
oc logs -n openshift-user-workload-monitoring prometheus-operator-7c7b77cbff-7v7km -c prometheus-operator

# check labels

oc get service prometheus-example-app -n servicemonitor-a --show-labels

# follow up

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