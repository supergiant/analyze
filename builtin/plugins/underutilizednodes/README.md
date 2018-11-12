#Underutilized nodes sunsetting plugin

Works only when ProviderID is set.
it matches EC2 instances with k8s cluster minions based on next equality:
k8s node -> ProviderID === ec2 instance -> Instance ID