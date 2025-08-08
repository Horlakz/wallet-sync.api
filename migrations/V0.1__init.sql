-- Strict mode for good measure
SET
    sql_mode = 'STRICT_ALL_TABLES';

-- Drop existing tables if any (for development reset)
DROP TABLE IF EXISTS reconciliation_logs;

DROP TABLE IF EXISTS ledger_entries;

DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS accounts;

DROP TABLE IF EXISTS users;

-- Users Table
CREATE TABLE
    users (
        id CHAR(36) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL
    );

-- Accounts Table
CREATE TABLE
    accounts (
        id CHAR(36) PRIMARY KEY,
        user_id CHAR(36) NOT NULL,
        account_type ENUM ('wallet', 'fee', 'reserve') DEFAULT 'wallet' NOT NULL,
        currency VARCHAR(10) DEFAULT 'NGN' NOT NULL,
        balance DECIMAL(32, 2) NOT NULL DEFAULT 0.00,
        number VARCHAR(20) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- Transactions Table
CREATE TABLE
    transactions (
        id CHAR(36) PRIMARY KEY,
        user_id CHAR(36) NOT NULL,
        reference VARCHAR(255) UNIQUE NOT NULL,
        type ENUM ('credit', 'debit') NOT NULL,
        status ENUM ('pending', 'completed', 'failed') DEFAULT 'pending' NOT NULL,
        amount DECIMAL(32, 2) NOT NULL,
        currency VARCHAR(10) DEFAULT 'NGN' NOT NULL,
        description VARCHAR(255) NOT NULL,
        metadata JSON,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- Ledger Entries Table
CREATE TABLE
    ledger_entries (
        id CHAR(36) PRIMARY KEY,
        user_id CHAR(36) NOT NULL,
        account_id CHAR(36) NOT NULL,
        transaction_id CHAR(36) NOT NULL,
        entry_type ENUM ('debit', 'credit') NOT NULL,
        amount DECIMAL(32, 2) NOT NULL,
        description VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (account_id) REFERENCES accounts (id),
        FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

-- Reconciliation Logs Table
CREATE TABLE
    reconciliation_logs (
        id CHAR(36) PRIMARY KEY,
        account_id CHAR(36) NOT NULL,
        computed_balance DECIMAL(32, 2) NOT NULL,
        stored_balance DECIMAL(32, 2) NOT NULL,
        discrepancy DECIMAL(32, 2) GENERATED ALWAYS AS (stored_balance - computed_balance) STORED,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL,
        FOREIGN KEY (account_id) REFERENCES accounts (id)
    );