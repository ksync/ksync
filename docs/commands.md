- `ksync init`

    - Sets the cluster up by starting the radar daemonset.
    - Starts `ksync watch` in the background. TODO

- `ksync create`

    Add a spec to sync. This gets watched and started/stopped automatically.

- `ksync delete`

    Remove a spec to sync.

- `ksync get`

    Fetch the status of all current specs

- `ksync watch`

    Watch for matching pods in the background (based off pod name and selector). Start syncs for any that come online.

- `ksync doctor`

    Debug what's happening under the covers and look for any possible issues with the system.

- `ksync version`

    Print out version information for the local binary. If the server binary is reachable and healthy, print information for that as well.
