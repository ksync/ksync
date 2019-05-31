# Running ksync in an RBAC enabled cluster

## Localhost
RBAC roles and bindings to create to make ksync work from localhost

- Create a `ksync` clusterrole

```
cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ksync
rules:
- apiGroups:
  - extensions
  resources:
  - daemonsets
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - list
- apiGroups:
  - ""
  resources:
  - pods/portforward
  verbs:
  - create
  EOF
 ```

- Please change the user name in the below `yaml` and bind the developer to the `ksync` clusterrole

```
cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ksync
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ksync
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: developer
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: useremail@domain
EOF
```

## Remote

To install the `ksync` remote components, cluster admin role is required to perform the `ksync` init.
```
cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-admin
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
EOF
```

Please change the user name in the below `yaml` and bind the developer to the `cluster-admin` role
```
cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ksync
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: developer
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: useremail@domain
EOF
```
