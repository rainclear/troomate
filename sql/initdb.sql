CREATE TABLE Currencies (
  id INTEGER PRIMARY KEY,
  Currency varchar(4) UNIQUE NOT NULL
);

CREATE TABLE AccountTypes (
  id INTEGER PRIMARY KEY,
  AccountType varchar(256) UNIQUE NOT NULL
);

CREATE TABLE AccountRoles (
  id INTEGER PRIMARY KEY,
  AccountRole varchar(256) UNIQUE NOT NULL
);

CREATE TABLE Accounts (
  id INTEGER PRIMARY KEY,
  AccountName varchar(256) UNIQUE NOT NULL,
  Balance decimal(20, 2) NOT NULL DEFAULT 0,
  OpeningDate date NOT NULL,
  OpeningBalance decimal(20, 2) NOT NULL DEFAULT 0,
  AccountNumber varchar(256),
  Currency_id int NOT NULL,
  AccountType_id int NOT NULL,
  AccountRole_id int NOT NULL,
  Notes varchar(512),
  FOREIGN KEY (Currency_id) REFERENCES Currencies(id),
  FOREIGN KEY (AccountType_id) REFERENCES AccountTypes(id),
  FOREIGN KEY (AccountRole_id) REFERENCES AccountRoles(id)
);

CREATE TABLE TransactionCategories (
  id INTEGER PRIMARY KEY,
  TransCategory varchar UNIQUE NOT NULL
);

CREATE TABLE TransactionTypes (
  id INTEGER PRIMARY KEY,
  TransType varchar UNIQUE NOT NULL
);

CREATE TABLE Transactions (
  id INTEGER PRIMARY KEY,
  Description varchar(256),
  TransDate date NOT NULL,
  Amount decimal(20, 2) NOT NULL DEFAULT 0,
  FromAccount_id bigint NOT NULL,
  ToAccount_id bigint NOT NULL CHECK(ToAccount_id <> FromAccount_id),
  TransCategory_id int NOT NULL,
  TransType_id int NOT NULL,
  Notes varchar(512),
  FOREIGN KEY (FromAccount_id) REFERENCES Accounts(id),
  FOREIGN KEY (ToAccount_id) REFERENCES Accounts(id),
  FOREIGN KEY (TransCategory_id) REFERENCES TransactionCategories(id),
  FOREIGN KEY (TransType_id) REFERENCES TransactionTypes(id)
);

CREATE INDEX ACCTNAME ON Accounts (AccountName);

CREATE INDEX FROMACCT ON Transactions (FromAccount_id);

CREATE INDEX TOACCT ON Transactions (ToAccount_id);

CREATE INDEX ACCTPAIR ON Transactions (FromAccount_id, ToAccount_id);

INSERT INTO Currencies (Currency) VALUES ('CAD');

INSERT INTO Currencies (Currency) VALUES ('USD');

INSERT INTO AccountTypes (AccountType) VALUES ('Assets');

INSERT INTO AccountTypes (AccountType) VALUES ('Equity');

INSERT INTO AccountTypes (AccountType) VALUES ('Expenses');

INSERT INTO AccountTypes (AccountType) VALUES ('Imbalance');

INSERT INTO AccountTypes (AccountType) VALUES ('Income');

INSERT INTO AccountTypes (AccountType) VALUES ('Liabilities');

INSERT INTO AccountRoles (AccountRole) VALUES ('Chequing');

INSERT INTO AccountRoles (AccountRole) VALUES ('Savings');

INSERT INTO AccountRoles (AccountRole) VALUES ('TFSA');

INSERT INTO TransactionTypes (TransType) VALUES ('Expense');

INSERT INTO TransactionTypes (TransType) VALUES ('Income');

INSERT INTO TransactionTypes (TransType) VALUES ('Transfer');

INSERT INTO TransactionCategories (TransCategory) VALUES ("Grocery");

INSERT INTO TransactionCategories (TransCategory) VALUES ("Utilities");