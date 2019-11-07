Jenkins Operator on Openshift

Due to the way the Openshift specific image for Jenkins works, we have to create a few extra resources and define some new resources as well.

### Jenkins on Openshift Custom Resource

Below, we have defined the CR for Jenkins which works with the Jenkins image provided by the Openshift Team. This image has some preconfigured plugins and defined mechanisms layered on top of the Jenkins LTS, due to which the extra environment variables and such are needed.

```
apiVersion: jenkins.io/v1alpha2
kind: Jenkins
metadata:
  generation: 1
  name: jenkins-openshift
spec:
  backup:
    action: {}
    containerName: ""
    interval: 0
    makeBackupBeforePodDeletion: false
  configurationAsCode:
    configurations: null
    secret:
      name: ""
  groovyScripts:
    configurations: null
    secret:
      name: ""
  master:
    basePlugins:
    - name: kubernetes
      version: 1.18.3
    - name: workflow-job
      version: "2.34"
    - name: workflow-aggregator
      version: "2.6"
    - name: git
      version: 3.12.0
    - name: job-dsl
      version: "1.76"
    - name: configuration-as-code
      version: "1.29"
    - name: configuration-as-code-support
      version: "1.19"
    - name: kubernetes-credentials-provider
      version: 0.12.1
    containers:
    - command:
      - /usr/bin/go-init
      - -main
      - /usr/libexec/s2i/run
      env:
      - name: OPENSHIFT_ENABLE_OAUTH
        value: "true"
      - name: OPENSHIFT_ENABLE_REDIRECT_PROMPT
        value: "true"
      - name: DISABLE_ADMINISTRATIVE_MONITORS        
        value: "false"
      - name: KUBERNETES_MASTER
        value: "https://kubernetes.default:443"
      - name: KUBERNETES_TRUST_CERTIFICATES
        value: "true"
      - name: JENKINS_SERVICE_NAME
        value: "jenkins-operator-http-jenkins-openshift"
      - name: JNLP_SERVICE_NAME
        value: "jenkins-operator-slave-jenkins-openshift"
      - name: JENKINS_UC_INSECURE
        value: "false"
      - name: JENKINS_HOME
        value: "/var/lib/jenkins"
      - name: JAVA_OPTS
        value: -XX:+UnlockExperimentalVMOptions -XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:MaxRAMFraction=1 -Djenkins.install.runSetupWizard=false -Djava.awt.headless=true
      image: quay.io/openshift/origin-jenkins:latest
      imagePullPolicy: Always
      livenessProbe:
        httpGet:
          path: /login
          port: 8080
          scheme: HTTP
        initialDelaySeconds: 420
        periodSeconds: 360
        timeoutSeconds: 240
      name: jenkins-master
      readinessProbe:
        httpGet:
          path: /login
          port: 8080
          scheme: HTTP
        initialDelaySeconds: 3
        periodSeconds: 0
        timeoutSeconds: 240
      resources:
        limits:
          cpu: 600m
          memory: 4Gi
        requests:
          cpu: 500m
          memory: 3Gi
    securityContext:
      runAsUser: 1000700001
  restore:
    action: {}
    containerName: ""
  service:
    port: 8080
    type: ClusterIP
  slaveService:
    port: 50000
    type: ClusterIP
```

Do refer to https://docs.openshift.com/container-platform/4.1/openshift_images/managing-images/using-image-pull-secrets.html for configuring Image Pull Secrets.

In Openshift the default SCC set is "restricted" which only allows UIDs to be set > 1000000000 in which case we cannot set the user for a Pod to be root. This can be set accordingly by the user based on their requirements.

### _Rolebinding_ to allow use of openshift specific plugins in Jenkins

The following _RoleBinding_ will bind the `edit` _Role_ to the _ServiceAccount_ of the Jenkins instance which has been created by the Operator. 

First we need to find the name of the _ServiceAccount_ you would ike to bind the `edit` _Role_ to.
```
$ oc get sa
NAME                                 SECRETS   AGE
builder                              2         142m
default                              2         142m
deployer                             2         142m
jenkins-operator-jenkins-openshift   2         8m12s
```
The ServiceAccount named `jenkins-operator-jenkins-openshift` is what we are looking for.

```yaml
oc create -f - << EOF
apiVersion: authorization.openshift.io/v1
kind: RoleBinding
metadata:
  name: jenkins_edit
roleRef:
  name: edit
subjects:
- kind: ServiceAccount
  name: jenkins-operator-jenkins-openshift
EOF
```

Use the above command to create the _RoleBinding_.

### Changes in _ServiceAccount_ and _Route_ to make Openshift Login Plugin work smoothly

To make the [openshift-login-plugin](https://plugins.jenkins.io/openshift-login) work smoothly with Jenkins we have to edit the Created _ServiceAccount_ and also create a _Route_ for Redirection and TLS Edge Termination without any issues.

```
oc annotate serviceaccount jenkins-operator-jenkins-openshift "serviceaccounts.openshift.io/oauth-redirectreference.jenkins"="{\"kind\":\"OAuthRedirectReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"Route\",\"name\":\"jenkins-operator-jenkins-openshift\"}}"
```

Next we create the _Route_.
First we need to figure out for which Service we are creating a Route for. SO run the following command and pick out the http service whcih matches your instance of Jenkins created by the Operator. 
```
$ oc get svc
NAME                                         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE                                                   
jenkins-operator-http-jenkins-openshift    ClusterIP   172.30.53.12   <none>        8080/TCP    28m                                                   
jenkins-operator-slave-jenkins-openshift   ClusterIP   172.30.41.13   <none>        50000/TCP   28m    
```

Let's pick out the _Service_ with the https in the name and create the _ROute accordingly.


```yaml
apiVersion: route.openshift.io/v1
kind: Route
  labels:
    app: jenkins-operator
    jenkins-cr: jenkins-openshift
  name: jenkins-operator-http-jenkins-openshift
  namespace: jenkins-operator-test
spec:
  tls:
    insecureEdgeTerminationPolicy: "Redirect"
    termination: "edge"
  port:
    targetPort: 8080
  to:
    kind: Service
    name: jenkins-operator-http-jenkins-openshift
EOF
```

Use the following command to create the _Route_.

After the creation of the _Route_. It can be used to navigate to the Jenkins Login Page and login with your Openshift Credentials.
