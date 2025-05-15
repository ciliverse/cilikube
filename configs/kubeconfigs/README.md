# Kubeconfigs

This directory stores the kubeconfig files uploaded by users for different Kubernetes clusters managed by CiliKube.
Each file is typically named after the cluster it represents.
Do not manually edit files in this directory unless you know what you are doing.

# Kubeconfig File Structure 
The kubeconfig file is a YAML file that contains the following sections:
- `apiVersion`: The version of the Kubernetes API that the kubeconfig file is using.
- `clusters`: A list of clusters that the kubeconfig file can connect to. Each cluster entry contains:
  - `cluster`: The name of the cluster.
  - `server`: The URL of the Kubernetes API server.
  - `certificate-authority-data`: The base64-encoded CA certificate for the cluster.
- `contexts`: A list of contexts that the kubeconfig file can use. Each context entry contains:
    - `context`: The name of the context.
    - `cluster`: The name of the cluster that this context is associated with.
    - `user`: The name of the user that this context is associated with.        
- `current-context`: The name of the current context that is being used.
- `kind`: The kind of the object. For kubeconfig files, this is usually "Config".
- `preferences`: User preferences for the kubeconfig file.
- `users`: A list of users that the kubeconfig file can use. Each user entry contains:
  - `name`: The name of the user.
  - `user`: The authentication information for the user, which can include:
    - `client-certificate-data`: The base64-encoded client certificate for the user.
    - `client-key-data`: The base64-encoded client key for the user.
    - `token`: A bearer token for authentication.
    - `username`: The username for basic authentication.
    - `password`: The password for basic authentication.
- `extensions`: A list of extensions that the kubeconfig file can use. 