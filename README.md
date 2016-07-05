MViews
==

Sets up MViews, a materialized views generator, refresher, and API
endpoint for postgres and redshift.

I found the need for this project while working with Looker endpoints
and not being able to rely on their data responses inside of internal
applications.

This projects mimics Looker Endpoint behavior and tries to be a drop in
replacement for Looker Endpoints, except that we don't drop tables while
trying to get a response from that very table.

The logic flows as such:

- Refresh tables on a timer
- When starting refresh, rebuild the staging table
- Once staging table is rebuilt, move pointer to read from the staging
table.
- Copy staging table to primary table, then shift reading pointer back
to the primary.

Results are cached internally in the application, which isn't feasible
in larger endpoint deployments but works nicely for testing. In
production deployment, this should be placed behind a caching layer.

LICENSE: MIT (See LICENSE file)
