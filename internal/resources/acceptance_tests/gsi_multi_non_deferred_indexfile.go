package acceptance_tests

var GsiMultiNonDeferredIndexfile = `

    
        resource "couchbase-capella_query_indexes" "idx_non_deferred1" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index1"
          index_keys      = ["c1"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred2" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index2"
          index_keys      = ["c2"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred3" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index3"
          index_keys      = ["c3"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred4" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index4"
          index_keys      = ["c4"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred5" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index5"
          index_keys      = ["c5"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred6" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index6"
          index_keys      = ["c6"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred7" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index7"
          index_keys      = ["c7"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred8" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index8"
          index_keys      = ["c8"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred9" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index9"
          index_keys      = ["c9"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred10" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index10"
          index_keys      = ["c10"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred11" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index11"
          index_keys      = ["c11"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred12" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index12"
          index_keys      = ["c12"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred13" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index13"
          index_keys      = ["c13"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred14" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index14"
          index_keys      = ["c14"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred15" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index15"
          index_keys      = ["c15"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred16" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index16"
          index_keys      = ["c16"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred17" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index17"
          index_keys      = ["c17"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred18" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index18"
          index_keys      = ["c18"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred19" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index19"
          index_keys      = ["c19"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred20" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index20"
          index_keys      = ["c20"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred21" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index21"
          index_keys      = ["c21"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred22" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index22"
          index_keys      = ["c22"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred23" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index23"
          index_keys      = ["c23"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred24" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index24"
          index_keys      = ["c24"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred25" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index25"
          index_keys      = ["c25"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred26" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index26"
          index_keys      = ["c26"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred27" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index27"
          index_keys      = ["c27"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred28" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index28"
          index_keys      = ["c28"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred29" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index29"
          index_keys      = ["c29"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred30" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index30"
          index_keys      = ["c30"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred31" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index31"
          index_keys      = ["c31"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred32" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index32"
          index_keys      = ["c32"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred33" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index33"
          index_keys      = ["c33"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred34" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index34"
          index_keys      = ["c34"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred35" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index35"
          index_keys      = ["c35"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred36" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index36"
          index_keys      = ["c36"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred37" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index37"
          index_keys      = ["c37"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred38" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index38"
          index_keys      = ["c38"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred39" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index39"
          index_keys      = ["c39"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred40" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index40"
          index_keys      = ["c40"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred41" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index41"
          index_keys      = ["c41"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred42" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index42"
          index_keys      = ["c42"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred43" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index43"
          index_keys      = ["c43"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred44" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index44"
          index_keys      = ["c44"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred45" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index45"
          index_keys      = ["c45"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred46" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index46"
          index_keys      = ["c46"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred47" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index47"
          index_keys      = ["c47"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred48" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index48"
          index_keys      = ["c48"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred49" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index49"
          index_keys      = ["c49"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred50" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index50"
          index_keys      = ["c50"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred51" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index51"
          index_keys      = ["c51"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred52" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index52"
          index_keys      = ["c52"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred53" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index53"
          index_keys      = ["c53"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred54" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index54"
          index_keys      = ["c54"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred55" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index55"
          index_keys      = ["c55"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred56" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index56"
          index_keys      = ["c56"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred57" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index57"
          index_keys      = ["c57"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred58" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index58"
          index_keys      = ["c58"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred59" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index59"
          index_keys      = ["c59"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred60" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index60"
          index_keys      = ["c60"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred61" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index61"
          index_keys      = ["c61"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred62" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index62"
          index_keys      = ["c62"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred63" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index63"
          index_keys      = ["c63"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred64" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index64"
          index_keys      = ["c64"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred65" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index65"
          index_keys      = ["c65"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred66" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index66"
          index_keys      = ["c66"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred67" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index67"
          index_keys      = ["c67"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred68" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index68"
          index_keys      = ["c68"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred69" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index69"
          index_keys      = ["c69"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred70" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index70"
          index_keys      = ["c70"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred71" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index71"
          index_keys      = ["c71"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred72" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index72"
          index_keys      = ["c72"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred73" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index73"
          index_keys      = ["c73"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred74" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index74"
          index_keys      = ["c74"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred75" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index75"
          index_keys      = ["c75"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred76" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index76"
          index_keys      = ["c76"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred77" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index77"
          index_keys      = ["c77"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred78" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index78"
          index_keys      = ["c78"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred79" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index79"
          index_keys      = ["c79"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred80" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index80"
          index_keys      = ["c80"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred81" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index81"
          index_keys      = ["c81"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred82" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index82"
          index_keys      = ["c82"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred83" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index83"
          index_keys      = ["c83"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred84" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index84"
          index_keys      = ["c84"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred85" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index85"
          index_keys      = ["c85"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred86" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index86"
          index_keys      = ["c86"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred87" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index87"
          index_keys      = ["c87"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred88" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index88"
          index_keys      = ["c88"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred89" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index89"
          index_keys      = ["c89"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred90" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index90"
          index_keys      = ["c90"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred91" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index91"
          index_keys      = ["c91"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred92" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index92"
          index_keys      = ["c92"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred93" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index93"
          index_keys      = ["c93"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred94" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index94"
          index_keys      = ["c94"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred95" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index95"
          index_keys      = ["c95"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred96" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index96"
          index_keys      = ["c96"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred97" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index97"
          index_keys      = ["c97"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred98" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index98"
          index_keys      = ["c98"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred99" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index99"
          index_keys      = ["c99"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
        
        resource "couchbase-capella_query_indexes" "idx_non_deferred100" {
          organization_id = var.organization_id
          project_id      = couchbase-capella_project.terraform_project.id
          cluster_id      = couchbase-capella_cluster.new_cluster.id
          bucket_name     = "travel-sample"
          scope_name      = "inventory"
          collection_name = "airline"
          index_name      = "non_deferred_index100"
          index_keys      = ["c100"]
          partition_by    = ["meta().id"]
          where           = "geo.alt > 1000"
          with = {
            num_partition = 3
          }
        }
`
