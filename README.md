# go-ipfs-ds-azure
Go implementation of ipfs Azure datastore

## Prerequisites
- Golang https://golang.org/dl/
- Make https://www.gnu.org/software/make/


## How to use it?
1. download go-ipfs source code from https://github.com/ipfs/go-ipfs
2. cd into go-ipfs
3. use command `go get github.com/decentralized-identity/go-ipfs-ds-azure@latest`
4. add `azureds github.com/decentralized-identity/go-ipfs-ds-azure/plugin 0` to the following file "go-ipfs/plugin/loader/preload_list"
5. cd back to go-ipfs
6. use command `make install`
7. ipfs should have the azureds plugin now

## How to configure it?
Once ipfs is successfully built, use the following configurations

1. add the following to store_spec in ipfs
`{"mounts":[{"accountName":"INSERT_ACCOUNT_NAME","containerName":"INSERT_CONTAINER_NAME","folderName":"INSERT_FOLDER_NAME","mountpoint":"/blocks"},{"mountpoint":"/","path":"datastore","type":"levelds"}],"type":"mount"}`

2. replace the default block mount point with the following in config for ipfs under Spec.mounts
`        {
          "child": {
            "accountName": "INSERT_ACCOUNT_NAME",
            "accountKey": "INSERT_ACCOUNT_KEY",
            "containerName": "INSERT_CONTAINER_NAME",
            "folderName": "INSERT_FOLDER_NAME",
            "type": "azure-data-store"
          },
          "mountpoint": "/blocks",
          "prefix": "azure.datastore",
          "type": "measure"
        },`

3. (optional) it is also possible to switch the root mount point from levelds to Azure, but block is the bare minimum for data to be persistent.

4. now you can run ipfs with the config file and it will connect to Azure `ipfs init CONFIG_FILE_NAME`
