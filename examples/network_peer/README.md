# Capella Network Peers Example AWS

This example shows how to create and manage Network Peers in Capella for AWS.

This creates a new network peer in the selected Capella AWS cluster and lists existing network peers in the cluster. It uses the cluster id to create and list network peers.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new network peer in Capella as stated in the `create_network_peer.tf` file.
2. UPDATE: Update the network peer configuration using Terraform.
3. LIST: List existing network peer in Capella as stated in the `list_network_peers.tf` file.
4. IMPORT: Import a network peer that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created network peer from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.1.0.0/23"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00ff"
            }
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.1.0.0/23"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00ff"
            }
          + gcp_config = null
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)


─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new network peer

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.1.0.0/23"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00ff"
            }
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.1.0.0/23"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00ff"
            }
          + gcp_config = null
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Creation complete after 6s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-06-26 18:13:42.685598255 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-26 18:13:49.087263425 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
    "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestAWS"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = {
      "account_id" = "123456789123"
      "cidr" = "10.1.0.0/23"
      "provider_id" = "pcx-12345678912345678"
      "region" = "us-east-1"
      "vpc_id" = "vpc-141f0fffff141aa00ff"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "aws"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-000000000000"

```

### Note the peer id for the new Network Peer
Command: `terraform output new_network_peer`

Sample Output:
```
$ terraform output new_network_peer
{
  "audit" = {
    "created_at" = "2024-06-26 18:13:42.685598255 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-26 18:13:49.087263425 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
    "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestAWS"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = {
      "account_id" = "123456789123"
      "cidr" = "10.1.0.0/23"
      "provider_id" = "pcx-12345678912345678"
      "region" = "us-east-1"
      "vpc_id" = "vpc-141f0fffff141aa00ff"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "aws"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}

```

In this case, the peer ID for my new network peer is `ffffffff-aaaa-1414-eeee-000000000000`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_network_peers.existing_network_peers
couchbase-capella_network_peer.new_network_peer

```

## IMPORT
### Remove the resource `new_network_peer` from the Terraform State file

Command: `terraform state rm couchbase-capella_network_peer.new_network_peer`

Sample Output:
```
$ terraform state rm couchbase-capella_network_peer.new_network_peer
Removed couchbase-capella_network_peer.new_network_peer
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_network_peer.new_network_peer id=<peer_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$ terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_network_peer.new_network_peer: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_network_peer.new_network_peer: Import prepared!
  Prepared couchbase-capella_network_peer for import
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the peer ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the network peer belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

```

## UPDATE
### Let us edit the terraform.tfvars file to change the network peer configuration settings. 

### Changed the vpc_id and cidr, destroys and forces replacement on terraform apply.

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer must be replaced
-/+ resource "couchbase-capella_network_peer" "new_network_peer" {
      ~ audit           = {
          ~ created_at  = "2024-06-26 18:13:42.685598255 +0000 UTC" -> (known after apply)
          ~ created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> (known after apply)
          ~ modified_at = "2024-06-26 18:13:49.087263425 +0000 UTC" -> (known after apply)
          ~ modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> (known after apply)
          ~ version     = 2 -> (known after apply)
        } -> (known after apply)
      ~ commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
        ] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
        name            = "VPCPeerTFTestAWS"
      ~ provider_config = { # forces replacement
          ~ aws_config = {
              ~ cidr        = "10.1.0.0/23" -> "10.2.0.0/23"
              ~ provider_id = "pcx-12345678912345678" -> (known after apply)
              ~ vpc_id      = "vpc-141f0fffff141aa00ff" -> "vpc-12345678912345678"
                # (2 unchanged attributes hidden)
            }
        }
      ~ status          = {
          ~ reasoning = "" -> (known after apply)
          ~ state     = "complete" -> (known after apply)
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_network_peer   = {
      ~ audit           = {
          - created_at  = "2024-06-26 18:13:42.685598255 +0000 UTC"
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - modified_at = "2024-06-26 18:13:49.087263425 +0000 UTC"
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - version     = 2
        } -> (known after apply)
      ~ commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
        ] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
        name            = "VPCPeerTFTestAWS"
      ~ provider_config = {
          ~ aws_config = {
              ~ cidr        = "10.1.0.0/23" -> "10.2.0.0/23"
              ~ provider_id = "pcx-12345678912345678" -> (known after apply)
              ~ vpc_id      = "vpc-141f0fffff141aa00ff" -> "vpc-12345678912345678"
                # (2 unchanged attributes hidden)
            }
            # (1 unchanged attribute hidden)
        }
      ~ status          = {
          - reasoning = ""
          - state     = "complete"
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }
  ~ peer_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 5s
couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Creation complete after 5s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2024-06-26 18:13:42.685598255 +0000 UTC"
        "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
        "modified_at" = "2024-06-26 18:13:49.087263425 +0000 UTC"
        "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
        "version" = 2
      }
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "name" = "VPCPeerTFTestAWS"
      "provider_config" = {
        "aws_config" = {
          "account_id" = "123456789123"
          "cidr" = "10.1.0.0/23"
          "provider_id" = "pcx-12345678912345678"
          "region" = "us-east-1"
          "vpc_id" = "vpc-141f0fffff141aa00ff"
        }
        "gcp_config" = null /* object */
      }
      "provider_type" = ""
      "status" = {
        "reasoning" = ""
        "state" = "complete"
      }
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-06-26 18:39:20.745080509 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-26 18:39:25.532316011 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-1234567891234567",
    "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-12345678912345678,VPCRegion=us-east-1 --region=us-east-1",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestAWS"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = {
      "account_id" = "123456789123"
      "cidr" = "10.2.0.0/23"
      "provider_id" = "pcx-1234567891234567"
      "region" = "us-east-1"
      "vpc_id" = "vpc-12345678912345678"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "aws"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-000000000000"


```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be destroyed
  - resource "couchbase-capella_network_peer" "new_network_peer" {
      - audit           = {
          - created_at  = "2024-06-26 18:39:20.745080509 +0000 UTC" -> null
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - modified_at = "2024-06-26 18:39:25.532316011 +0000 UTC" -> null
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - version     = 2 -> null
        } -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-1234567891234567",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-12345678912345678,VPCRegion=us-east-1 --region=us-east-1",
        ] -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "VPCPeerTFTestAWS" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - provider_config = {
          - aws_config = {
              - account_id  = "123456789123" -> null
              - cidr        = "10.2.0.0/23" -> null
              - provider_id = "pcx-1234567891234567" -> null
              - region      = "us-east-1" -> null
              - vpc_id      = "vpc-12345678912345678" -> null
            } -> null
        } -> null
      - provider_type   = "aws" -> null
      - status          = {
          - reasoning = "" -> null
          - state     = "complete" -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - network_peers_list = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2024-06-26 18:39:20.745080509 +0000 UTC"
                  - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - modified_at = "2024-06-26 18:39:25.532316011 +0000 UTC"
                  - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - version     = 2
                }
              - id              = "ffffffff-aaaa-1414-eeee-000000000000"
              - name            = "VPCPeerTFTestAWS"
              - provider_config = {
                  - aws_config = {
                      - account_id  = "123456789123"
                      - cidr        = "10.2.0.0/23"
                      - provider_id = "pcx-1234567891234567"
                      - region      = "us-east-1"
                      - vpc_id      = "vpc-12345678912345678"
                    }
                  - gcp_config = null
                }
              - provider_type   = ""
              - status          = {
                  - reasoning = ""
                  - state     = "complete"
                }
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null
  - new_network_peer   = {
      - audit           = {
          - created_at  = "2024-06-26 18:39:20.745080509 +0000 UTC"
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - modified_at = "2024-06-26 18:39:25.532316011 +0000 UTC"
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - version     = 2
        }
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-1234567891234567",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-12345678912345678,VPCRegion=us-east-1 --region=us-east-1",
        ]
      - id              = "ffffffff-aaaa-1414-eeee-000000000000"
      - name            = "VPCPeerTFTestAWS"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - provider_config = {
          - aws_config = {
              - account_id  = "123456789123"
              - cidr        = "10.2.0.0/23"
              - provider_id = "pcx-1234567891234567"
              - region      = "us-east-1"
              - vpc_id      = "vpc-12345678912345678"
            }
          - gcp_config = null
        }
      - provider_type   = "aws"
      - status          = {
          - reasoning = ""
          - state     = "complete"
        }
    } -> null
  - peer_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 4s

Destroy complete! Resources: 1 destroyed.

```

# Capella Network Peers Example GCP

This example shows how to create and manage Network Peers in Capella for GCP.

This creates a new network peer in the selected Capella GCP cluster and lists existing network peers in the cluster. It uses the cluster id to create and list network peers.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new network peer in Capella as stated in the `create_network_peer.tf` file.
2. UPDATE: Update the network peer configuration using Terraform.
3. LIST: List existing network peer in Capella as stated in the `list_network_peers.tf` file.
4. IMPORT: Import a network peer that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created network peer from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan   
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestGCP"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + gcp_config = {
              + cidr            = "10.0.4.0/23"
              + network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
              + project_id      = "test-project-id"
              + provider_id     = (known after apply)
              + service_account = "service-account-name@project-id.iam.gserviceaccount.com"
            }
        }
      + provider_type   = "gcp"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestGCP"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = null
          + gcp_config = {
              + cidr            = "10.0.4.0/23"
              + network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
              + project_id      = "test-project-id"
              + provider_id     = (known after apply)
              + service_account = "service-account-name@project-id.iam.gserviceaccount.com"
            }
        }
      + provider_type   = "gcp"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

### Apply the Plan, in order to create a new network peer

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestGCP"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + gcp_config = {
              + cidr            = "10.0.4.0/23"
              + network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
              + project_id      = "test-project-id"
              + provider_id     = (known after apply)
              + service_account = "service-account-name@project-id.iam.gserviceaccount.com"
            }
        }
      + provider_type   = "gcp"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestGCP"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws_config = null
          + gcp_config = {
              + cidr            = "10.0.4.0/23"
              + network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
              + project_id      = "test-project-id"
              + provider_id     = (known after apply)
              + service_account = "service-account-name@project-id.iam.gserviceaccount.com"
            }
        }
      + provider_type   = "gcp"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Still creating... [10s elapsed]
couchbase-capella_network_peer.new_network_peer: Creation complete after 19s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-06-27 19:34:00.202552084 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-27 19:34:18.661844385 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-f5be7fe9-f5a0-4468 --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-ffffffff-aaaa-1414-eeee-000000000000",
    "gcloud dns managed-zones create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-project=test-project-id --visibility=private",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestGCP"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "gcp_config" = {
      "cidr" = "10.0.4.0/23"
      "network_name" = "cc-ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "test-project-id"
      "provider_id" = "cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414"
      "service_account" = "service-account-name@project-id.iam.gserviceaccount.com"
    }
  }
  "provider_type" = "gcp"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-000000000000"

```

### Note the peer id for the new Network Peer
Command: `terraform output new_network_peer`

Sample Output:
```
$ terraform output new_network_peer
{
  "audit" = {
    "created_at" = "2024-06-27 19:34:00.202552084 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-27 19:34:18.661844385 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-f5be7fe9-f5a0-4468 --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-ffffffff-aaaa-1414-eeee-000000000000",
    "gcloud dns managed-zones create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-project=test-project-id --visibility=private",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestGCP"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "gcp_config" = {
      "cidr" = "10.0.4.0/23"
      "network_name" = "cc-ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "test-project-id"
      "provider_id" = "cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414"
      "service_account" = "service-account-name@project-id.iam.gserviceaccount.com"
    }
  }
  "provider_type" = "gcp"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}

```

In this case, the peer ID for my new network peer is `ffffffff-aaaa-1414-eeee-000000000000`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_network_peers.existing_network_peers
couchbase-capella_network_peer.new_network_peer

```

## IMPORT
### Remove the resource `new_network_peer` from the Terraform State file

Command: `terraform state rm couchbase-capella_network_peer.new_network_peer`

Sample Output:
```
$ terraform state rm couchbase-capella_network_peer.new_network_peer
Removed couchbase-capella_network_peer.new_network_peer
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_network_peer.new_network_peer id=<peer_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$  terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000

couchbase-capella_network_peer.new_network_peer: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Import prepared!
  Prepared couchbase-capella_network_peer for import
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the peer ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the network peer belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

```

## UPDATE
### Let us edit the terraform.tfvars file to change the network peer configuration settings.

# Changed the cluster_id, destroys and forces replacement

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer must be replaced
-/+ resource "couchbase-capella_network_peer" "new_network_peer" {
      ~ audit           = {
          ~ created_at  = "2024-06-29 00:22:24.928834346 +0000 UTC" -> (known after apply)
          ~ created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> (known after apply)
          ~ modified_at = "2024-06-29 00:22:42.698786674 +0000 UTC" -> (known after apply)
          ~ modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> (known after apply)
          ~ version     = 2 -> (known after apply)
        } -> (known after apply)
      ~ cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> "562a53c8-fe49-444f-8230-dc976fa1f749" # forces replacement
      ~ commands        = [
          - "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-f5be7fe9-f5a0-4468 --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-ffffffff-aaaa-1414-eeee-000000000000",
          - "gcloud dns managed-zones create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-project=test-project-id --visibility=private",
        ] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
        name            = "VPCPeerTFTestGCP"
      ~ provider_config = { # forces replacement
          ~ gcp_config = {
              ~ provider_id     = "cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414" -> (known after apply)
                # (4 unchanged attributes hidden)
            }
        }
      ~ status          = {
          ~ reasoning = "" -> (known after apply)
          ~ state     = "complete" -> (known after apply)
        } -> (known after apply)
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ network_peers_list = {
      ~ cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> "562a53c8-fe49-444f-8230-dc976fa1f749"
      ~ data            = [
          - {
              - audit           = {
                  - created_at  = "2024-06-29 00:22:24.928834346 +0000 UTC"
                  - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - modified_at = "2024-06-29 00:22:42.698786674 +0000 UTC"
                  - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - version     = 2
                }
              - id              = "ffffffff-aaaa-1414-eeee-000000000000"
              - name            = "VPCPeerTFTestGCP"
              - provider_config = {
                  - aws_config = null
                  - gcp_config = {
                      - cidr            = "10.0.4.0/23"
                      - network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
                      - project_id      = "test-project-id"
                      - provider_id     = "cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414"
                      - service_account = "service-account-name@project-id.iam.gserviceaccount.com"
                    }
                }
              - status          = {
                  - reasoning = ""
                  - state     = "complete"
                }
            },
        ] -> null
        # (2 unchanged attributes hidden)
    }
  ~ new_network_peer   = {
      ~ audit           = {
          - created_at  = "2024-06-29 00:22:24.928834346 +0000 UTC"
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - modified_at = "2024-06-29 00:22:42.698786674 +0000 UTC"
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - version     = 2
        } -> (known after apply)
      ~ cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> "562a53c8-fe49-444f-8230-dc976fa1f749"
      ~ commands        = [
          - "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-f5be7fe9-f5a0-4468 --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-ffffffff-aaaa-1414-eeee-000000000000",
          - "gcloud dns managed-zones create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-project=test-project-id --visibility=private",
        ] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
        name            = "VPCPeerTFTestGCP"
      ~ provider_config = {
          ~ gcp_config = {
              ~ provider_id     = "cc-ffffffff-aaaa-1414-eeee-000000000000-cc-ffffffff-aaaa-1414" -> (known after apply)
                # (4 unchanged attributes hidden)
            }
            # (1 unchanged attribute hidden)
        }
      ~ status          = {
          - reasoning = ""
          - state     = "complete"
        } -> (known after apply)
        # (3 unchanged attributes hidden)
    }
  ~ peer_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 10s elapsed]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 13s
couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Still creating... [10s elapsed]
couchbase-capella_network_peer.new_network_peer: Creation complete after 20s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "562a53c8-fe49-444f-8230-dc976fa1f749"
  "data" = tolist(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-06-29 00:37:45.4338168 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-29 00:38:04.635673378 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "562a53c8-fe49-444f-8230-dc976fa1f749"
  "commands" = toset([
    "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-562a53c8-fe49-444f --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-562a53c8-fe49-444f-8230-dc976fa1f749",
    "gcloud dns managed-zones create cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-562a53c8-fe49-444f-8230-dc976fa1f749 --target-project=test-project-id --visibility=private",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestGCP"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "gcp_config" = {
      "cidr" = "10.0.4.0/23"
      "network_name" = "cc-ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "test-project-id"
      "provider_id" = "cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414"
      "service_account" = "service-account-name@project-id.iam.gserviceaccount.com"
    }
  }
  "provider_type" = "gcp"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-000000000000"

```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be destroyed
  - resource "couchbase-capella_network_peer" "new_network_peer" {
      - audit           = {
          - created_at  = "2024-06-29 00:37:45.4338168 +0000 UTC" -> null
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - modified_at = "2024-06-29 00:38:04.635673378 +0000 UTC" -> null
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - version     = 2 -> null
        } -> null
      - cluster_id      = "562a53c8-fe49-444f-8230-dc976fa1f749" -> null
      - commands        = [
          - "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-562a53c8-fe49-444f --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-562a53c8-fe49-444f-8230-dc976fa1f749",
          - "gcloud dns managed-zones create cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-562a53c8-fe49-444f-8230-dc976fa1f749 --target-project=test-project-id --visibility=private",
        ] -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "VPCPeerTFTestGCP" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - provider_config = {
          - gcp_config = {
              - cidr            = "10.0.4.0/23" -> null
              - network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000" -> null
              - project_id      = "test-project-id" -> null
              - provider_id     = "cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414" -> null
              - service_account = "service-account-name@project-id.iam.gserviceaccount.com" -> null
            } -> null
        } -> null
      - provider_type   = "gcp" -> null
      - status          = {
          - reasoning = "" -> null
          - state     = "complete" -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - network_peers_list = {
      - cluster_id      = "562a53c8-fe49-444f-8230-dc976fa1f749"
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2024-06-29 00:37:45.4338168 +0000 UTC"
                  - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - modified_at = "2024-06-29 00:38:04.635673378 +0000 UTC"
                  - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
                  - version     = 2
                }
              - id              = "ffffffff-aaaa-1414-eeee-000000000000"
              - name            = "VPCPeerTFTestGCP"
              - provider_config = {
                  - aws_config = null
                  - gcp_config = {
                      - cidr            = "10.0.4.0/23"
                      - network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
                      - project_id      = "test-project-id"
                      - provider_id     = "cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414"
                      - service_account = "service-account-name@project-id.iam.gserviceaccount.com"
                    }
                }
              - status          = {
                  - reasoning = ""
                  - state     = "complete"
                }
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null
  - new_network_peer   = {
      - audit           = {
          - created_at  = "2024-06-29 00:37:45.4338168 +0000 UTC"
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - modified_at = "2024-06-29 00:38:04.635673378 +0000 UTC"
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - version     = 2
        }
      - cluster_id      = "562a53c8-fe49-444f-8230-dc976fa1f749"
      - commands        = [
          - "gcloud compute networks peerings create cc-ffffffff-aaaa-1414-eeee-000000000000-cc-562a53c8-fe49-444f --network=cc-ffffffff-aaaa-1414-eeee-000000000000 --peer-project test-project-id --peer-network cc-562a53c8-fe49-444f-8230-dc976fa1f749",
          - "gcloud dns managed-zones create cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414 --description=\"Peering Zone to Capella\" --dns-name=test-dns-name --account=service-account-name@project-id.iam.gserviceaccount.com --networks=cc-ffffffff-aaaa-1414-eeee-000000000000 --target-network=cc-562a53c8-fe49-444f-8230-dc976fa1f749 --target-project=test-project-id --visibility=private",
        ]
      - id              = "ffffffff-aaaa-1414-eeee-000000000000"
      - name            = "VPCPeerTFTestGCP"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - provider_config = {
          - aws_config = null
          - gcp_config = {
              - cidr            = "10.0.4.0/23"
              - network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
              - project_id      = "test-project-id"
              - provider_id     = "cc-562a53c8-fe49-444f-8230-dc976fa1f749-cc-ffffffff-aaaa-1414"
              - service_account = "service-account-name@project-id.iam.gserviceaccount.com"
            }
        }
      - provider_type   = "gcp"
      - status          = {
          - reasoning = ""
          - state     = "complete"
        }
    } -> null
  - peer_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 10s elapsed]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 14s

Destroy complete! Resources: 1 destroyed.

```

# Capella Network Peers Example Azure

This example shows how to create and manage Network Peers in Capella for Azure.

This creates a new network peer in the selected Capella Azure cluster and lists existing network peers in the cluster. It uses the cluster id to create and list network peers.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new network peer in Capella as stated in the `create_network_peer.tf` file.
2. UPDATE: Update the network peer configuration using Terraform.
3. LIST: List existing network peer in Capella as stated in the `list_network_peers.tf` file.
4. IMPORT: Import a network peer that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created network peer from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace//go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VNETPeerTFTestAZURE"
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + provider_config = {
          + azure_config = {
              + cidr            = "10.6.0.0/16"
              + provider_id     = (known after apply)
              + resource_group  = "test-rg"
              + subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
              + tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
              + vnet_id         = "test-vnet"
            }
        }
      + provider_type   = "azure"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2024-08-23 23:30:39.790107885 +0000 UTC"
                  + created_by  = "sample_apikey"
                  + modified_at = "2024-08-23 23:35:39.359925093 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
                  + version     = 6
                }
              + id              = "ffffffff-aaaa-1414-eeee-00000000000"
              + name            = "VNETPeerTFTestAZURE"
              + provider_config = {
                  + aws_config   = null
                  + azure_config = {
                      + cidr            = "10.6.0.0/16"
                      + provider_id     = ""
                      + resource_group  = "\"test-rg\""
                      + subscription_id = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      + tenant_id       = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      + vnet_id         = "\"test-vnet\""
                    }
                  + gcp_config   = null
                }
              + status          = {
                  + reasoning = ""
                  + state     = "failed"
                }
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VNETPeerTFTestAZURE"
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + provider_config = {
          + aws_config   = null
          + azure_config = {
              + cidr            = "10.6.0.0/16"
              + provider_id     = (known after apply)
              + resource_group  = "test-rg"
              + subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
              + tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
              + vnet_id         = "test-vnet"
            }
          + gcp_config   = null
        }
      + provider_type   = "azure"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

### Apply the Plan, in order to create a new network peer

Command: `terraform apply`

Sample Output:
```
$  terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VNETPeerTFTestAZURE"
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + provider_config = {
          + azure_config = {
              + cidr            = "10.6.0.0/16"
              + provider_id     = (known after apply)
              + resource_group  = "test-rg"
              + subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
              + tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
              + vnet_id         = "test-vnet"
            }
        }
      + provider_type   = "azure"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + network_peers_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2024-08-23 23:30:39.790107885 +0000 UTC"
                  + created_by  = "sample_apikey"
                  + modified_at = "2024-08-23 23:35:39.359925093 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
                  + version     = 6
                }
              + id              = "ffffffff-aaaa-1414-eeee-00000000000"
              + name            = "VNETPeerTFTestAZURE"
              + provider_config = {
                  + aws_config   = null
                  + azure_config = {
                      + cidr            = "10.6.0.0/16"
                      + provider_id     = ""
                      + resource_group  = "\"test-rg\""
                      + subscription_id = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      + tenant_id       = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      + vnet_id         = "\"test-vnet\""
                    }
                  + gcp_config   = null
                }
              + status          = {
                  + reasoning = ""
                  + state     = "failed"
                }
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
    }
  + new_network_peer   = {
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VNETPeerTFTestAZURE"
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      + provider_config = {
          + aws_config   = null
          + azure_config = {
              + cidr            = "10.6.0.0/16"
              + provider_id     = (known after apply)
              + resource_group  = "test-rg"
              + subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
              + tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
              + vnet_id         = "test-vnet"
            }
          + gcp_config   = null
        }
      + provider_type   = "azure"
      + status          = (known after apply)
    }
  + peer_id            = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Still creating... [10s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [20s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [30s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [40s elapsed]
couchbase-capella_network_peer.new_network_peer: Creation complete after 43s [id=ffffffff-aaaa-1414-eeee-00000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2024-08-23 23:30:39.790107885 +0000 UTC"
        "created_by" = "sample_apikey"
        "modified_at" = "2024-08-23 23:35:39.359925093 +0000 UTC"
        "modified_by" = "ffffffff-aaaa-1414-eeee-00000000000"
        "version" = 6
      }
      "id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "name" = "VNETPeerTFTestAZURE"
      "provider_config" = {
        "aws_config" = null /* object */
        "azure_config" = {
          "cidr" = "10.6.0.0/16"
          "provider_id" = ""
          "resource_group" = "\"test-rg\""
          "subscription_id" = "\"ffffffff-aaaa-1414-eeee-00000000000\""
          "tenant_id" = "\"ffffffff-aaaa-1414-eeee-00000000000\""
          "vnet_id" = "\"test-vnet\""
        }
        "gcp_config" = null /* object */
      }
      "status" = {
        "reasoning" = ""
        "state" = "failed"
      }
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-08-24 00:01:23.524305127 +0000 UTC"
    "created_by" = "sample_apikey"
    "modified_at" = "2024-08-24 00:02:05.606678591 +0000 UTC"
    "modified_by" = "sample_apikey"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "commands" = toset([])
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "name" = "VNETPeerTFTestAZURE"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "azure_config" = {
      "cidr" = "10.6.0.0/16"
      "provider_id" = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
      "resource_group" = "test-rg"
      "subscription_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "tenant_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "vnet_id" = "test-vnet"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "azure"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-00000000000"

```

### Note the peer id for the new Network Peer
Command: `terraform output new_network_peer`

Sample Output:
```
$ terraform output new_network_peer
{
  "audit" = {
    "created_at" = "2024-08-24 00:01:23.524305127 +0000 UTC"
    "created_by" = "sample_apikey"
    "modified_at" = "2024-08-24 00:02:05.606678591 +0000 UTC"
    "modified_by" = "sample_apikey"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "commands" = toset([])
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "name" = "VNETPeerTFTestAZURE"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "azure_config" = {
      "cidr" = "10.6.0.0/16"
      "provider_id" = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
      "resource_group" = "test-rg"
      "subscription_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "tenant_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "vnet_id" = "test-vnet"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "azure"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}

```

In this case, the peer ID for my new network peer is `ffffffff-aaaa-1414-eeee-000000000000`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_network_peers.existing_network_peers
couchbase-capella_network_peer.new_network_peer

```

## IMPORT
### Remove the resource `new_network_peer` from the Terraform State file

Command: `terraform state rm couchbase-capella_network_peer.new_network_peer`

Sample Output:
```
$ terraform state rm couchbase-capella_network_peer.new_network_peer
Removed couchbase-capella_network_peer.new_network_peer
Successfully removed 1 resource instance(s).

```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_network_peer.new_network_peer id=<peer_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$  terraform import couchbase-capella_network_peer.new_network_peer id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000

couchbase-capella_network_peer.new_network_peer: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Import prepared!
  Prepared couchbase-capella_network_peer for import
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the peer ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the network peer belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

```

## UPDATE
### Let us edit the terraform.tfvars file to change the network peer configuration settings.

# Changed the network_name, destroys and forces replacement

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-00000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer must be replaced
-/+ resource "couchbase-capella_network_peer" "new_network_peer" {
      ~ audit           = {
          ~ created_at  = "2024-08-24 00:01:23.524305127 +0000 UTC" -> (known after apply)
          ~ created_by  = "sample_apikey" -> (known after apply)
          ~ modified_at = "2024-08-24 00:02:05.606678591 +0000 UTC" -> (known after apply)
          ~ modified_by = "sample_apikey" -> (known after apply)
          ~ version     = 2 -> (known after apply)
        } -> (known after apply)
      ~ commands        = [] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-00000000000" -> (known after apply)
      ~ name            = "VNETPeerTFTestAZURE" -> "VNETPeerTFTestAZURE2" # forces replacement
      ~ provider_config = { # forces replacement
          ~ azure_config = {
              ~ provider_id     = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet" -> (known after apply)
                # (5 unchanged attributes hidden)
            }
        }
      ~ status          = {
          ~ reasoning = "" -> (known after apply)
          ~ state     = "complete" -> (known after apply)
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_network_peer   = {
      ~ audit           = {
          - created_at  = "2024-08-24 00:01:23.524305127 +0000 UTC"
          - created_by  = "sample_apikey"
          - modified_at = "2024-08-24 00:02:05.606678591 +0000 UTC"
          - modified_by = "sample_apikey"
          - version     = 2
        } -> (known after apply)
      ~ commands        = [] -> (known after apply)
      ~ id              = "ffffffff-aaaa-1414-eeee-00000000000" -> (known after apply)
      ~ name            = "VNETPeerTFTestAZURE" -> "VNETPeerTFTestAZURE2"
      ~ provider_config = {
          ~ azure_config = {
              ~ provider_id     = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet" -> (known after apply)
                # (5 unchanged attributes hidden)
            }
            # (2 unchanged attributes hidden)
        }
      ~ status          = {
          - reasoning = ""
          - state     = "complete"
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }
  ~ peer_id            = "ffffffff-aaaa-1414-eeee-00000000000" -> (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-00000000000]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 10s elapsed]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 20s elapsed]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 26s
couchbase-capella_network_peer.new_network_peer: Creating...
couchbase-capella_network_peer.new_network_peer: Still creating... [10s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [20s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [30s elapsed]
couchbase-capella_network_peer.new_network_peer: Still creating... [40s elapsed]
couchbase-capella_network_peer.new_network_peer: Creation complete after 42s [id=ffffffff-aaaa-1414-eeee-00000000000]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

network_peers_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2024-08-23 23:30:39.790107885 +0000 UTC"
        "created_by" = "sample_apikey"
        "modified_at" = "2024-08-23 23:35:39.359925093 +0000 UTC"
        "modified_by" = "418653cc-a3c5-4943-9ff0-ee4d36a49ca3"
        "version" = 6
      }
      "id" = "c3319197-22ad-4513-b56d-b5b660c749f5"
      "name" = "VNETPeerTFTestAZURE"
      "provider_config" = {
        "aws_config" = null /* object */
        "azure_config" = {
          "cidr" = "10.6.0.0/16"
          "provider_id" = ""
          "resource_group" = "\"test-rg\""
          "subscription_id" = "\"ffffffff-aaaa-1414-eeee-00000000000\""
          "tenant_id" = "\"ffffffff-aaaa-1414-eeee-00000000000\""
          "vnet_id" = "\"test-vnet\""
        }
        "gcp_config" = null /* object */
      }
      "status" = {
        "reasoning" = ""
        "state" = "failed"
      }
    },
    {
      "audit" = {
        "created_at" = "2024-08-24 00:01:23.524305127 +0000 UTC"
        "created_by" = "sample_apikey"
        "modified_at" = "2024-08-24 00:02:05.606678591 +0000 UTC"
        "modified_by" = "sample_apikey"
        "version" = 2
      }
      "id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "name" = "VNETPeerTFTestAZURE"
      "provider_config" = {
        "aws_config" = null /* object */
        "azure_config" = {
          "cidr" = "10.6.0.0/16"
          "provider_id" = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
          "resource_group" = "test-rg"
          "subscription_id" = "ffffffff-aaaa-1414-eeee-00000000000"
          "tenant_id" = "ffffffff-aaaa-1414-eeee-00000000000"
          "vnet_id" = "test-vnet"
        }
        "gcp_config" = null /* object */
      }
      "status" = {
        "reasoning" = ""
        "state" = "complete"
      }
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
}
new_network_peer = {
  "audit" = {
    "created_at" = "2024-08-24 00:10:56.415670004 +0000 UTC"
    "created_by" = "sample_apikey"
    "modified_at" = "2024-08-24 00:11:37.523742051 +0000 UTC"
    "modified_by" = "sample_apikey"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "commands" = toset([])
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "name" = "VNETPeerTFTestAZURE2"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "provider_config" = {
    "aws_config" = null /* object */
    "azure_config" = {
      "cidr" = "10.6.0.0/16"
      "provider_id" = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
      "resource_group" = "test-rg"
      "subscription_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "tenant_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "vnet_id" = "test-vnet"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "azure"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-00000000000"


```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$  terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_network_peers.existing_network_peers: Reading...
couchbase-capella_network_peer.new_network_peer: Refreshing state... [id=ffffffff-aaaa-1414-eeee-00000000000]
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be destroyed
  - resource "couchbase-capella_network_peer" "new_network_peer" {
      - audit           = {
          - created_at  = "2024-08-24 00:10:56.415670004 +0000 UTC" -> null
          - created_by  = "sample_apikey" -> null
          - modified_at = "2024-08-24 00:11:37.523742051 +0000 UTC" -> null
          - modified_by = "sample_apikey" -> null
          - version     = 2 -> null
        } -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - commands        = [] -> null
      - id              = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - name            = "VNETPeerTFTestAZURE2" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - provider_config = {
          - azure_config = {
              - cidr            = "10.6.0.0/16" -> null
              - provider_id     = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet" -> null
              - resource_group  = "test-rg" -> null
              - subscription_id = "ffffffff-aaaa-1414-eeee-00000000000" -> null
              - tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000" -> null
              - vnet_id         = "test-vnet" -> null
            } -> null
        } -> null
      - provider_type   = "azure" -> null
      - status          = {
          - reasoning = "" -> null
          - state     = "complete" -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - network_peers_list = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2024-08-24 00:10:56.415670004 +0000 UTC"
                  - created_by  = "sample_apikey"
                  - modified_at = "2024-08-24 00:11:37.523742051 +0000 UTC"
                  - modified_by = "sample_apikey"
                  - version     = 2
                }
              - id              = "ffffffff-aaaa-1414-eeee-00000000000"
              - name            = "VNETPeerTFTestAZURE2"
              - provider_config = {
                  - aws_config   = null
                  - azure_config = {
                      - cidr            = "10.6.0.0/16"
                      - provider_id     = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
                      - resource_group  = "test-rg"
                      - subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
                      - tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
                      - vnet_id         = "test-vnet"
                    }
                  - gcp_config   = null
                }
              - status          = {
                  - reasoning = ""
                  - state     = "complete"
                }
            },
          - {
              - audit           = {
                  - created_at  = "2024-08-23 23:30:39.790107885 +0000 UTC"
                  - created_by  = "sample_apikey"
                  - modified_at = "2024-08-23 23:35:39.359925093 +0000 UTC"
                  - modified_by = "418653cc-a3c5-4943-9ff0-ee4d36a49ca3"
                  - version     = 6
                }
              - id              = "c3319197-22ad-4513-b56d-b5b660c749f5"
              - name            = "VNETPeerTFTestAZURE"
              - provider_config = {
                  - aws_config   = null
                  - azure_config = {
                      - cidr            = "10.6.0.0/16"
                      - provider_id     = ""
                      - resource_group  = "\"test-rg\""
                      - subscription_id = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      - tenant_id       = "\"ffffffff-aaaa-1414-eeee-00000000000\""
                      - vnet_id         = "\"test-vnet\""
                    }
                  - gcp_config   = null
                }
              - status          = {
                  - reasoning = ""
                  - state     = "failed"
                }
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
    } -> null
  - new_network_peer   = {
      - audit           = {
          - created_at  = "2024-08-24 00:10:56.415670004 +0000 UTC"
          - created_by  = "sample_apikey"
          - modified_at = "2024-08-24 00:11:37.523742051 +0000 UTC"
          - modified_by = "sample_apikey"
          - version     = 2
        }
      - cluster_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      - commands        = []
      - id              = "ffffffff-aaaa-1414-eeee-00000000000"
      - name            = "VNETPeerTFTestAZURE2"
      - organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
      - provider_config = {
          - aws_config   = null
          - azure_config = {
              - cidr            = "10.6.0.0/16"
              - provider_id     = "/subscriptions/ffffffff-aaaa-1414-eeee-00000000000/resourceGroups/rg-ffffffff-aaaa-1414-eeee-00000000000/providers/Microsoft.Network/virtualNetworks/cc-ffffffff-aaaa-1414-eeee-00000000000/virtualNetworkPeerings/cc-ffffffff-aaaa-1414-eeee-00000000000-test-vnet"
              - resource_group  = "test-rg"
              - subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
              - tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
              - vnet_id         = "test-vnet"
            }
          - gcp_config   = null
        }
      - provider_type   = "azure"
      - status          = {
          - reasoning = ""
          - state     = "complete"
        }
    } -> null
  - peer_id            = "ffffffff-aaaa-1414-eeee-00000000000" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_network_peer.new_network_peer: Destroying... [id=ffffffff-aaaa-1414-eeee-00000000000]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 10s elapsed]
couchbase-capella_network_peer.new_network_peer: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 20s elapsed]
couchbase-capella_network_peer.new_network_peer: Destruction complete after 26s

Destroy complete! Resources: 1 destroyed.


```
