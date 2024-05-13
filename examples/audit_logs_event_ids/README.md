# Capella Clusters Audit Log Event IDs Example

This example shows how to retrive audit log event ids for a Capella cluster.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the audit log event ids for the cluster.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## GET

Command: `terraform plan`

Sample Output:
```

terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_event_ids.existing_auditlogeventids: Reading...
data.couchbase-capella_audit_log_event_ids.existing_auditlogeventids: Read complete after 1s

Changes to Outputs:
  + existing_auditlogeventids = {
      + cluster_id      = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + data            = [
          + {
              + description = "A N1QL ADVISE statement was executed"
              + id          = 28719
              + module      = "n1ql"
              + name        = "ADVISE statement"
            },
          + {
              + description = "A N1QL ALTER COLLECTION statement was executed"
              + id          = 36880
              + module      = "analytics"
              + name        = "ALTER COLLECTION statement"
            },
          + {
              + description = "A N1QL ALTER INDEX statement was executed"
              + id          = 28683
              + module      = "n1ql"
              + name        = "ALTER INDEX statement"
            },
          + {
              + description = "A N1QL BUILD INDEX statement was executed"
              + id          = 28684
              + module      = "n1ql"
              + name        = "BUILD INDEX statement"
            },
          + {
              + description = "A N1QL COMMIT TRANSACTION statement was executed"
              + id          = 28721
              + module      = "n1ql"
              + name        = "COMMIT TRANSACTION statement"
            },
          + {
              + description = "A N1QL CONNECT LINK statement was executed"
              + id          = 36877
              + module      = "analytics"
              + name        = "CONNECT LINK statement"
            },
          + {
              + description = "A N1QL CREATE COLLECTION statement was executed"
              + id          = 28715
              + module      = "n1ql"
              + name        = "CREATE COLLECTION statement"
            },
          + {
              + description = "A N1QL CREATE DATASET statement was executed"
              + id          = 36870
              + module      = "analytics"
              + name        = "CREATE DATASET statement"
            },
          + {
              + description = "A N1QL CREATE DATAVERSE statement was executed"
              + id          = 36868
              + module      = "analytics"
              + name        = "CREATE DATAVERSE statement"
            },
          + {
              + description = "A N1QL CREATE FUNCTION statement was executed"
              + id          = 28706
              + module      = "n1ql"
              + name        = "CREATE FUNCTION statement"
            },
          + {
              + description = "A N1QL CREATE INDEX statement was executed"
              + id          = 28681
              + module      = "n1ql"
              + name        = "CREATE INDEX statement"
            },
          + {
              + description = "A N1QL CREATE INDEX statement was executed"
              + id          = 36872
              + module      = "analytics"
              + name        = "CREATE INDEX statement"
            },
          + {
              + description = "A N1QL CREATE PRIMARY INDEX statement was executed"
              + id          = 28688
              + module      = "n1ql"
              + name        = "CREATE PRIMARY INDEX statement"
            },
          + {
              + description = "A N1QL CREATE SCOPE statement was executed"
              + id          = 28713
              + module      = "n1ql"
              + name        = "CREATE SCOPE statement"
            },
          + {
              + description = "A N1QL DELETE statement was executed"
              + id          = 28678
              + module      = "n1ql"
              + name        = "DELETE statement"
            },
          + {
              + description = "A N1QL DISCONNECT LINK statement was executed"
              + id          = 36878
              + module      = "analytics"
              + name        = "DISCONNECT LINK statement"
            },
          + {
              + description = "A N1QL DROP COLLECTION statement was executed"
              + id          = 28716
              + module      = "n1ql"
              + name        = "DROP COLLECTION statement"
            },
          + {
              + description = "A N1QL DROP DATASET statement was executed"
              + id          = 36871
              + module      = "analytics"
              + name        = "DROP DATASET statement"
            },
          + {
              + description = "A N1QL DROP DATAVERSE statement was executed"
              + id          = 36869
              + module      = "analytics"
              + name        = "DROP DATAVERSE statement"
            },
          + {
              + description = "A N1QL DROP FUNCTION statement was executed"
              + id          = 28707
              + module      = "n1ql"
              + name        = "DROP FUNCTION statement"
            },
          + {
              + description = "A N1QL DROP INDEX statement was executed"
              + id          = 28682
              + module      = "n1ql"
              + name        = "DROP INDEX statement"
            },
          + {
              + description = "A N1QL DROP INDEX statement was executed"
              + id          = 36873
              + module      = "analytics"
              + name        = "DROP INDEX statement"
            },
          + {
              + description = "A N1QL DROP SCOPE statement was executed"
              + id          = 28714
              + module      = "n1ql"
              + name        = "DROP SCOPE statement"
            },
          + {
              + description = "A N1QL EXECUTE FUNCTION statement was executed"
              + id          = 28708
              + module      = "n1ql"
              + name        = "EXECUTE FUNCTION statement"
            },
          + {
              + description = "A N1QL EXPLAIN statement was executed"
              + id          = 28673
              + module      = "n1ql"
              + name        = "EXPLAIN statement"
            },
          + {
              + description = "A N1QL FLUSH COLLECTION statement was executed"
              + id          = 28717
              + module      = "n1ql"
              + name        = "FLUSH COLLECTION statement"
            },
          + {
              + description = "A N1QL GRANT ROLE statement was executed"
              + id          = 28685
              + module      = "n1ql"
              + name        = "GRANT ROLE statement"
            },
          + {
              + description = "A N1QL INFER statement was executed"
              + id          = 28675
              + module      = "n1ql"
              + name        = "INFER statement"
            },
          + {
              + description = "A N1QL INSERT statement was executed"
              + id          = 28676
              + module      = "n1ql"
              + name        = "INSERT statement"
            },
          + {
              + description = "A N1QL MERGE statement was executed"
              + id          = 28680
              + module      = "n1ql"
              + name        = "MERGE statement"
            },
          + {
              + description = "A N1QL PREPARE statement was executed"
              + id          = 28674
              + module      = "n1ql"
              + name        = "PREPARE statement"
            },
          + {
              + description = "A N1QL REVOKE ROLE statement was executed"
              + id          = 28686
              + module      = "n1ql"
              + name        = "REVOKE ROLE statement"
            },
          + {
              + description = "A N1QL ROLLBACK TRANSACTION TO SAVEPOINT statement was executed"
              + id          = 28723
              + module      = "n1ql"
              + name        = "ROLLBACK TRANSACTION TO SAVEPOINT statement"
            },
          + {
              + description = "A N1QL ROLLBACK TRANSACTION statement was executed"
              + id          = 28722
              + module      = "n1ql"
              + name        = "ROLLBACK TRANSACTION statement"
            },
          + {
              + description = "A N1QL SAVEPOINT statement was executed"
              + id          = 28725
              + module      = "n1ql"
              + name        = "SAVEPOINT statement"
            },
          + {
              + description = "A N1QL SELECT statement was executed"
              + id          = 28672
              + module      = "n1ql"
              + name        = "SELECT statement"
            },
          + {
              + description = "A N1QL SELECT statement was executed"
              + id          = 36867
              + module      = "analytics"
              + name        = "SELECT statement"
            },
          + {
              + description = "A N1QL SET TRANSACTION ISOLATION statement was executed"
              + id          = 28724
              + module      = "n1ql"
              + name        = "SET TRANSACTION ISOLATION statement"
            },
          + {
              + description = "A N1QL START TRANSACTION statement was executed"
              + id          = 28720
              + module      = "n1ql"
              + name        = "START TRANSACTION statement"
            },
          + {
              + description = "A N1QL UPDATE STATISTICS statement was executed"
              + id          = 28718
              + module      = "n1ql"
              + name        = "UPDATE STATISTICS statement"
            },
          + {
              + description = "A N1QL UPDATE statement was executed"
              + id          = 28679
              + module      = "n1ql"
              + name        = "UPDATE statement"
            },
          + {
              + description = "A N1QL UPSERT statement was executed"
              + id          = 28677
              + module      = "n1ql"
              + name        = "UPSERT statement"
            },
          + {
              + description = "A backup plan was removed"
              + id          = 45060
              + module      = "backup"
              + name        = "Delete plan"
            },
          + {
              + description = "A document was retrieved from the repository backups"
              + id          = 45071
              + module      = "backup"
              + name        = "Examine repository"
            },
          + {
              + description = "A manual backup was triggered on an active repository"
              + id          = 45068
              + module      = "backup"
              + name        = "Backup repository"
            },
          + {
              + description = "A manual merge was triggered on an active repository"
              + id          = 45069
              + module      = "backup"
              + name        = "Merge repository"
            },
          + {
              + description = "A new active backup repository was added"
              + id          = 45062
              + module      = "backup"
              + name        = "Add repository"
            },
          + {
              + description = "A new backup plan was added"
              + id          = 45058
              + module      = "backup"
              + name        = "Add plan"
            },
          + {
              + description = "A repository was deleted"
              + id          = 45072
              + module      = "backup"
              + name        = "Delete repository"
            },
          + {
              + description = "A repository was fetched"
              + id          = 45066
              + module      = "backup"
              + name        = "Fetch repository"
            },
          + {
              + description = "A successful node configuration change was made."
              + id          = 36866
              + module      = "analytics"
              + name        = "Node configuration change"
            },
          + {
              + description = "A successful service configuration change was made."
              + id          = 36865
              + module      = "analytics"
              + name        = "Service configuration change"
            },
          + {
              + description = "A user has been denied access to the REST API due to invalid permissions or credentials"
              + id          = 45074
              + module      = "backup"
              + name        = "Access denied"
            },
          + {
              + description = "Access denied to the REST API due to invalid permissions or credentials"
              + id          = 40966
              + module      = "view_engine"
              + name        = "Access denied"
            },
          + {
              + description = "An HTTP request was made to archive or restore N1QL metadata "
              + id          = 28728
              + module      = "n1ql"
              + name        = "N1QL backup / restore API request"
            },
          + {
              + description = "An HTTP request was made to initate graceful shutdown"
              + id          = 28729
              + module      = "n1ql"
              + name        = "/admin/shutdown API request"
            },
          + {
              + description = "An HTTP request was made to run an FFDC collection"
              + id          = 28731
              + module      = "n1ql"
              + name        = "/admin/ffdc API request"
            },
          + {
              + description = "An HTTP request was made to run garbage collection"
              + id          = 28730
              + module      = "n1ql"
              + name        = "/admin/gc API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/active_requests."
              + id          = 28692
              + module      = "n1ql"
              + name        = "/admin/active_requests API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/clusters."
              + id          = 28701
              + module      = "n1ql"
              + name        = "/admin/clusters API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/completed_requests."
              + id          = 28702
              + module      = "n1ql"
              + name        = "/admin/completed_requests API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/config."
              + id          = 28698
              + module      = "n1ql"
              + name        = "/admin/config API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/dictionary_cache."
              + id          = 28711
              + module      = "n1ql"
              + name        = "/admin/dictionary_cache API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/functions."
              + id          = 28704
              + module      = "n1ql"
              + name        = "/admin/functions API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/active_requests."
              + id          = 28694
              + module      = "n1ql"
              + name        = "/admin/indexes/active_requests API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/completed_requests."
              + id          = 28695
              + module      = "n1ql"
              + name        = "/admin/indexes/completed_requests API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/dictionary_cache."
              + id          = 28712
              + module      = "n1ql"
              + name        = "/admin/indexes/dictionary_cache API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/functions."
              + id          = 28705
              + module      = "n1ql"
              + name        = "/admin/indexes/functions API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/prepareds."
              + id          = 28693
              + module      = "n1ql"
              + name        = "/admin/indexes/prepareds API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/tasks."
              + id          = 28710
              + module      = "n1ql"
              + name        = "/admin/indexes/tasks API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/indexes/transactions."
              + id          = 28727
              + module      = "n1ql"
              + name        = "/admin/indexes/transactions API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/ping."
              + id          = 28697
              + module      = "n1ql"
              + name        = "/admin/ping API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/prepareds."
              + id          = 28691
              + module      = "n1ql"
              + name        = "/admin/prepareds API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/settings."
              + id          = 28700
              + module      = "n1ql"
              + name        = "/admin/settings API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/ssl_cert."
              + id          = 28699
              + module      = "n1ql"
              + name        = "/admin/ssl_cert API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/stats."
              + id          = 28689
              + module      = "n1ql"
              + name        = "/admin/stats API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/tasks."
              + id          = 28709
              + module      = "n1ql"
              + name        = "/admin/tasks API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/transactions."
              + id          = 28726
              + module      = "n1ql"
              + name        = "/admin/transactions API request"
            },
          + {
              + description = "An HTTP request was made to the API at /admin/vitals."
              + id          = 28690
              + module      = "n1ql"
              + name        = "/admin/vitals API request"
            },
          + {
              + description = "An UNRECOGNIZED N1QL statement was encountered"
              + id          = 36879
              + module      = "analytics"
              + name        = "UNRECOGNIZED statement"
            },
          + {
              + description = "An active repository backup was deleted"
              + id          = 45073
              + module      = "backup"
              + name        = "Delete backup"
            },
          + {
              + description = "An active repository was archived"
              + id          = 45063
              + module      = "backup"
              + name        = "Archive repository"
            },
          + {
              + description = "An active repository was paused"
              + id          = 45064
              + module      = "backup"
              + name        = "Pause repository"
            },
          + {
              + description = "An active repository was resumed"
              + id          = 45065
              + module      = "backup"
              + name        = "Resume repository"
            },
          + {
              + description = "An alert email was successfully sent"
              + id          = 8257
              + module      = "ns_server"
              + name        = "alert email sent"
            },
          + {
              + description = "An unrecognized statement was received by the N1QL query engine"
              + id          = 28687
              + module      = "n1ql"
              + name        = "UNRECOGNIZED statement"
            },
          + {
              + description = "Authentication to the cluster succeeded"
              + id          = 20485
              + module      = "memcached"
              + name        = "authentication succeeded"
            },
          + {
              + description = "Backup service configuration was modified"
              + id          = 45056
              + module      = "backup"
              + name        = "Modify configuration"
            },
          + {
              + description = "Backup service configuration was retrieved"
              + id          = 45057
              + module      = "backup"
              + name        = "Fetch configuration"
            },
          + {
              + description = "Design Doc Meta Data Query Request"
              + id          = 40962
              + module      = "view_engine"
              + name        = "Query DDoc Meta Data"
            },
          + {
              + description = "Design Doc is Created"
              + id          = 40960
              + module      = "view_engine"
              + name        = "Create Design Doc"
            },
          + {
              + description = "Design Doc is Deleted"
              + id          = 40961
              + module      = "view_engine"
              + name        = "Delete Design Doc"
            },
          + {
              + description = "Design Doc is Updated"
              + id          = 40964
              + module      = "view_engine"
              + name        = "Update Design Doc"
            },
          + {
              + description = "Document was deleted"
              + id          = 20491
              + module      = "memcached"
              + name        = "document delete"
            },
          + {
              + description = "Document was locked"
              + id          = 20489
              + module      = "memcached"
              + name        = "document locked"
            },
          + {
              + description = "Document was modified"
              + id          = 20490
              + module      = "memcached"
              + name        = "document modify"
            },
          + {
              + description = "Document was mutated via the REST API"
              + id          = 8243
              + module      = "ns_server"
              + name        = "mutate document"
            },
          + {
              + description = "Document was read via the REST API"
              + id          = 8255
              + module      = "ns_server"
              + name        = "read document"
            },
          + {
              + description = "Document was read"
              + id          = 20488
              + module      = "memcached"
              + name        = "document read"
            },
          + {
              + description = "Existing backup plan was modified"
              + id          = 45059
              + module      = "backup"
              + name        = "Modify plan"
            },
          + {
              + description = "External user flushed the content of a memcached bucket"
              + id          = 20482
              + module      = "memcached"
              + name        = "external memcached bucket flush"
            },
          + {
              + description = "Information about the structure and contents of the backup repository was fetched."
              + id          = 45070
              + module      = "backup"
              + name        = "Info repository"
            },
          + {
              + description = "One or more backup plans where fetched"
              + id          = 45061
              + module      = "backup"
              + name        = "Fetch plan"
            },
          + {
              + description = "RBAC information was retrieved"
              + id          = 8265
              + module      = "ns_server"
              + name        = "RBAC information retrieved"
            },
          + {
              + description = "Rejected an invalid packet"
              + id          = 20483
              + module      = "memcached"
              + name        = "invalid packet"
            },
          + {
              + description = "Request to backup one or more eventing functions"
              + id          = 32793
              + module      = "eventing"
              + name        = "Backup Functions"
            },
          + {
              + description = "Request to create or update eventing function definition"
              + id          = 32768
              + module      = "eventing"
              + name        = "Create Function"
            },
          + {
              + description = "Request to delete eventing function definition"
              + id          = 32769
              + module      = "eventing"
              + name        = "Delete Function"
            },
          + {
              + description = "Request to delete eventing function draft definitions"
              + id          = 32773
              + module      = "eventing"
              + name        = "Delete Drafts"
            },
          + {
              + description = "Request to deploy eventing function"
              + id          = 32789
              + module      = "eventing"
              + name        = "Deploy Function"
            },
          + {
              + description = "Request to execute eventing node related functions"
              + id          = 32801
              + module      = "eventing"
              + name        = "Eventing System Event"
            },
          + {
              + description = "Request to export all eventing functions"
              + id          = 32785
              + module      = "eventing"
              + name        = "Export Functions"
            },
          + {
              + description = "Request to fetch eventing cluster stats"
              + id          = 32799
              + module      = "eventing"
              + name        = "Eventing Cluster Stats"
            },
          + {
              + description = "Request to fetch eventing config"
              + id          = 32780
              + module      = "eventing"
              + name        = "Fetch Config"
            },
          + {
              + description = "Request to fetch eventing deployed functions list"
              + id          = 32771
              + module      = "eventing"
              + name        = "List Deployed"
            },
          + {
              + description = "Request to fetch eventing function definition"
              + id          = 32770
              + module      = "eventing"
              + name        = "Fetch Functions"
            },
          + {
              + description = "Request to fetch eventing function draft definitions"
              + id          = 32772
              + module      = "eventing"
              + name        = "Fetch Drafts"
            },
          + {
              + description = "Request to fetch eventing function settings"
              + id          = 32783
              + module      = "eventing"
              + name        = "Get Settings"
            },
          + {
              + description = "Request to fetch eventing function stats"
              + id          = 32798
              + module      = "eventing"
              + name        = "Fetch Stats"
            },
          + {
              + description = "Request to fetch eventing function status"
              + id          = 32796
              + module      = "eventing"
              + name        = "Function Status"
            },
          + {
              + description = "Request to fetch eventing functions"
              + id          = 32795
              + module      = "eventing"
              + name        = "List Function"
            },
          + {
              + description = "Request to fetch eventing running function list"
              + id          = 32786
              + module      = "eventing"
              + name        = "List Running"
            },
          + {
              + description = "Request to get user eventing permissions"
              + id          = 32802
              + module      = "eventing"
              + name        = "Get User Info"
            },
          + {
              + description = "Request to import one or more eventing functions"
              + id          = 32784
              + module      = "eventing"
              + name        = "Import Functions"
            },
          + {
              + description = "Request to pause eventing function"
              + id          = 32791
              + module      = "eventing"
              + name        = "Pause Function"
            },
          + {
              + description = "Request to reset eventing function stats"
              + id          = 32797
              + module      = "eventing"
              + name        = "Clear Stats"
            },
          + {
              + description = "Request to restore one or more eventing functions from a backup"
              + id          = 32794
              + module      = "eventing"
              + name        = "Restore Functions"
            },
          + {
              + description = "Request to resume eventing function"
              + id          = 32792
              + module      = "eventing"
              + name        = "Resume Function"
            },
          + {
              + description = "Request to save a draft definition"
              + id          = 32774
              + module      = "eventing"
              + name        = "Save Draft"
            },
          + {
              + description = "Request to save eventing config"
              + id          = 32781
              + module      = "eventing"
              + name        = "Save Config"
            },
          + {
              + description = "Request to save settings for an eventing function"
              + id          = 32779
              + module      = "eventing"
              + name        = "Set Settings"
            },
          + {
              + description = "Request to start eventing function debugger"
              + id          = 32775
              + module      = "eventing"
              + name        = "Start Debug"
            },
          + {
              + description = "Request to start tracing eventing function execution"
              + id          = 32777
              + module      = "eventing"
              + name        = "Start Tracing"
            },
          + {
              + description = "Request to stop eventing function debugger"
              + id          = 32776
              + module      = "eventing"
              + name        = "Stop Debug"
            },
          + {
              + description = "Request to stop tracing eventing function execution"
              + id          = 32778
              + module      = "eventing"
              + name        = "Stop Tracing"
            },
          + {
              + description = "Request to undeploy eventing function"
              + id          = 32790
              + module      = "eventing"
              + name        = "Undeploy Function"
            },
          + {
              + description = "Session to the cluster has terminated"
              + id          = 20493
              + module      = "memcached"
              + name        = "session terminated"
            },
          + {
              + description = "The given tenant was rate limited"
              + id          = 20494
              + module      = "memcached"
              + name        = "tenant rate limited"
            },
          + {
              + description = "The repository data was restored"
              + id          = 45067
              + module      = "backup"
              + name        = "Restore repository"
            },
          + {
              + description = "The specified bucket was selected"
              + id          = 20492
              + module      = "memcached"
              + name        = "select bucket"
            },
          + {
              + description = "View Query Request"
              + id          = 40963
              + module      = "view_engine"
              + name        = "View Query"
            },
          + {
              + description = "opened DCP connection"
              + id          = 20480
              + module      = "memcached"
              + name        = "opened DCP connection"
            },
        ]
      + organization_id = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id      = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply"
now.
```