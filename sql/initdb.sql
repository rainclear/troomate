-- Currently one owner per db, so only one onwer record in one profile
CREATE TABLE Owners       (
  id                      INTEGER PRIMARY KEY,
  OwnerName               VARCHAR(512) UNIQUE NOT NULL,
  OwnerBirthDate          DATE NOT NULL Default '1900-01-01',
  TFSAStartDate           DATE NOT NULL Default '2009-01-01',
  TFSAEligibleDate        DATE NOT NULL Default '2009-01-01',
  Notes                   TEXT
);

CREATE TABLE Accounts     (
  id                      INTEGER PRIMARY KEY,
  OwnerId                 INTEGER NOT NULL,
  AccountName             VARCHAR(512) UNIQUE NOT NULL,
  Institution             VARCHAR(512),
  AccountNumber           VARCHAR(512),
  AccountNameAtCRA        VARCHAR(512) UNIQUE NOT NULL,
  AccountType             VARCHAR(512), -- Savings, GIC, Investment, ...
  AccountPurpose          VARCHAR(512), -- Daily, Short Term, Long Term, Retirement, ...
  OpeningDate             DATE,
  ClosingDate             DATE,
  Notes                   TEXT,

  FOREIGN KEY (OwnerId)   REFERENCES Owners(id)
);

CREATE TABLE Transactions (
  id                      INTEGER PRIMARY KEY,
  TransDate               DATE NOT NULL,
  Amount decimal(20, 2)   NOT NULL DEFAULT 0,
  AccountId               INTEGER NOT NULL,
  TransType               TEXT CHECK( TransType IN ('Deposit','Withdraw') ) NOT NULL DEFAULT 'Deposit',
  Notes                   TEXT,

  FOREIGN KEY (AccountId) REFERENCES Accounts(id)
);

CREATE INDEX OWNER        ON Owners       (OwnerName);
CREATE INDEX ACCTNAME     ON Accounts     (AccountName);
CREATE INDEX TRANSDATE    ON Transactions (TransDate);
