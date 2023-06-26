# Manifest

A small program for monitoring changes to your cloud environment. This tool may be used for unauthorized asset detection.

## Limitations
### Scans Missing Resources
Manifest uses the AWS SDK ResourceExplorer2 Search endpoint to list all resources within an environment. This endpoint is hardcoded to return no more than 1000 results. If your environment contains a high number of resources, it's possible that resources may not be detected if they are truncated by the search endpoint. In order to reduce the likelyhood of this, Manifest calls the search endpoint for each region and resource type directly (e.g `arn region:us-east-1 resourcetype:s3:bucket`). By breaking out the search into more queries with greater specificity, this risk should be avoided for many (if not most) cloud environments.

### Cron Scheduling
Manifest uses cron expressions to specify the interval to run cloud envionment scans within. Be careful not to set this interval too low or you run the risk of running multiple scans simultaniously.
