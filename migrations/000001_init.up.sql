CREATE TABLE storage_files (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  name VARCHAR(256) NOT NULL,
  extension VARCHAR(256) NOT NULL,
  size BIGINT NOT NULL DEFAULT(0),
  createdAt TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deletedAt TIMESTAMP(3)
)