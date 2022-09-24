provider "opensearch" {
  endpoints = "https://192.168.122.165:9200,https://192.168.122.166:9200"
  ca_cert = "../shared/opensearch-ca.pem"
  client_cert = "../shared/opensearch-admin.pem"
  client_key = "../shared/opensearch-admin.key"
}