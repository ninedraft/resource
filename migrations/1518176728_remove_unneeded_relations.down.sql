ALTER TYPE RESOURCE_KIND ADD VALUE 'domain' AFTER 'intservice';

CREATE TABLE IF NOT EXISTS deployment_volume (
  depl_id UUID NOT NULL,
  vol_id UUID NOT NULL,

  FOREIGN KEY (depl_id) REFERENCES deployments (id) ON DELETE CASCADE,
  FOREIGN KEY (vol_id) REFERENCES volumes (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS namespace_volume (
  ns_id UUID NOT NULL,
  vol_id UUID NOT NULL,

  FOREIGN KEY (ns_id) REFERENCES namespaces (id) ON DELETE CASCADE,
  FOREIGN KEY (vol_id) REFERENCES volumes (id) ON DELETE CASCADE
);
