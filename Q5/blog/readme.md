# Consensus-Breaking Change Explanation

I changed the key used to store the post in the underlying key-value store

If nodes running an older version of the code attempt to access this post using the old key format, they won't be able to find the post, leading to a consensus issue.