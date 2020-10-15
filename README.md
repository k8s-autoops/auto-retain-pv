# auto-retain-pv

## Usage

Create namespace `autoops` and apply yaml resources as described below.

```yaml
# create serviceaccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: auto-retain-pv
  namespace: autoops
---
# create clusterrole
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: auto-retain-pv
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["list", "patch"]
---
# create clusterrolebinding
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: auto-retain-pv
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: auto-retain-pv
subjects:
  - kind: ServiceAccount
    name: auto-retain-pv
    namespace: autoops
---
# create cronjob
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: auto-retain-pv
  namespace: autoops
spec:
  schedule: "*/5 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccount: auto-retain-pv
          containers:
            - name: auto-retain-pv
              image: autoops/auto-retain-pv
          restartPolicy: OnFailure
```

## Credits

Guo Y.K., MIT License
