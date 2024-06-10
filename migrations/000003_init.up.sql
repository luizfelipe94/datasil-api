ALTER TABLE
  storage_files
ADD
  COLUMN companyId VARCHAR(36) NOT NULL;

ALTER TABLE
  storage_files
ADD
  CONSTRAINT fk_files_company FOREIGN KEY (companyId) REFERENCES companies(id);

ALTER TABLE
  storage_files
ADD
  COLUMN path VARCHAR(512) NOT NULL DEFAULT '/';