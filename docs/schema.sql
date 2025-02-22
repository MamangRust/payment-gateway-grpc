CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "users" (
    "user_id" serial PRIMARY KEY,
    "firstname" VARCHAR(100) NOT NULL,
    "lastname" varchar(100) NOT NULL,
    "email" varchar(100) UNIQUE NOT NULL,
    "password" varchar(100) NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_firstname ON users (firstname);

CREATE INDEX idx_users_lastname ON users (lastname);

CREATE INDEX idx_users_firstname_lastname ON users (firstname, lastname);

CREATE INDEX idx_users_created_at ON users (created_at);

CREATE TABLE "cards" (
    "card_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "card_number" VARCHAR(16) UNIQUE NOT NULL,
    "card_type" VARCHAR(50) NOT NULL,
    "expire_date" DATE NOT NULL,
    "cvv" VARCHAR(3) NOT NULL,
    "card_provider" VARCHAR(50) NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_cards_card_number ON cards (card_number);

CREATE INDEX idx_cards_user_id ON cards (user_id);

CREATE INDEX idx_cards_card_type ON cards (card_type);

CREATE INDEX idx_cards_expire_date ON cards (expire_date);

CREATE INDEX idx_cards_user_id_card_type ON cards (user_id, card_type);

CREATE TABLE "merchants" (
    "merchant_id" SERIAL PRIMARY KEY,
    "merchant_no" UUID NOT NULL DEFAULT gen_random_uuid (),
    "name" VARCHAR(255) NOT NULL,
    "api_key" VARCHAR(255) UNIQUE NOT NULL,
    "user_id" INT NOT NULL REFERENCES "users" (user_id),
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_merchants_api_key ON merchants (api_key);

CREATE INDEX idx_merchants_user_id ON merchants (user_id);

CREATE INDEX idx_merchants_status ON merchants (status);

CREATE INDEX idx_merchants_name ON merchants (name);

CREATE INDEX idx_merchants_user_id_status ON merchants (user_id, status);

CREATE TABLE "saldos" (
    "saldo_id" SERIAL PRIMARY KEY,
    "card_number" VARCHAR(16) NOT NULL REFERENCES "cards" ("card_number"),
    "total_balance" INT NOT NULL,
    "withdraw_amount" INT DEFAULT 0,
    "withdraw_time" TIMESTAMP DEFAULT current_timestamp,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_saldos_card_number ON saldos (card_number);

CREATE INDEX idx_saldos_withdraw_time ON saldos (withdraw_time);

CREATE INDEX idx_saldos_total_balance ON saldos (total_balance);

CREATE INDEX idx_saldos_card_number_withdraw_time ON saldos (card_number, withdraw_time);

CREATE TABLE "transactions" (
    "transaction_id" SERIAL PRIMARY KEY,
    "transaction_no" UUID NOT NULL DEFAULT gen_random_uuid (),
    "card_number" VARCHAR(16) NOT NULL REFERENCES "cards" ("card_number"),
    "amount" INT NOT NULL,
    "payment_method" VARCHAR(50) NOT NULL,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id"),
    "transaction_time" TIMESTAMP NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_transactions_card_number ON transactions (card_number);

CREATE INDEX idx_transactions_merchant_id ON transactions (merchant_id);

CREATE INDEX idx_transactions_transaction_time ON transactions (transaction_time);

CREATE INDEX idx_transactions_payment_method ON transactions (payment_method);

CREATE INDEX idx_transactions_card_number_transaction_time ON transactions (card_number, transaction_time);

CREATE TABLE transfers (
    "transfer_id" SERIAL PRIMARY KEY,
    "transfer_no" UUID NOT NULL DEFAULT gen_random_uuid (),
    "transfer_from" VARCHAR(16) NOT NULL REFERENCES cards ("card_number"),
    "transfer_to" VARCHAR(16) NOT NULL REFERENCES cards ("card_number"),
    "transfer_amount" INT NOT NULL,
    "transfer_time" TIMESTAMP NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_transfers_transfer_from ON transfers (transfer_from);

CREATE INDEX idx_transfers_transfer_to ON transfers (transfer_to);

CREATE INDEX idx_transfers_transfer_time ON transfers (transfer_time);

CREATE INDEX idx_transfers_transfer_amount ON transfers (transfer_amount);

CREATE INDEX idx_transfers_transfer_from_transfer_time ON transfers (transfer_from, transfer_time);

CREATE TABLE "topups" (
    "topup_id" SERIAL PRIMARY KEY,
    "topup_no" UUID NOT NULL DEFAULT gen_random_uuid (),
    "card_number" VARCHAR(16) NOT NULL REFERENCES "cards" ("card_number"),
    "topup_amount" INT NOT NULL,
    "topup_method" VARCHAR(50) NOT NULL,
    "topup_time" TIMESTAMP NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_topups_card_number ON topups (card_number);

CREATE INDEX idx_topups_topup_no ON topups (topup_no);

CREATE INDEX idx_topups_topup_time ON topups (topup_time);

CREATE INDEX idx_topups_topup_method ON topups (topup_method);

CREATE INDEX idx_topups_card_number_topup_time ON topups (card_number, topup_time);

CREATE TABLE "withdraws" (
    "withdraw_id" SERIAL PRIMARY KEY,
    "withdraw_no" UUID NOT NULL DEFAULT gen_random_uuid (),
    "card_number" VARCHAR(16) NOT NULL REFERENCES cards ("card_number"),
    "withdraw_amount" INT NOT NULL,
    "withdraw_time" TIMESTAMP NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_withdraws_card_number ON withdraws (card_number);

CREATE INDEX idx_withdraws_withdraw_time ON withdraws (withdraw_time);

CREATE INDEX idx_withdraws_withdraw_amount ON withdraws (withdraw_amount);

CREATE INDEX idx_withdraws_card_number_withdraw_time ON withdraws (card_number, withdraw_time);

CREATE TABLE IF NOT EXISTS "roles" (
    "role_id" SERIAL PRIMARY KEY,
    "role_name" VARCHAR(50) UNIQUE NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_roles_role_name ON roles (role_name);

CREATE INDEX idx_roles_created_at ON roles (created_at);

CREATE INDEX idx_roles_updated_at ON roles (updated_at);

CREATE TABLE IF NOT EXISTS "user_roles" (
    "user_role_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id") ON DELETE CASCADE,
    "role_id" INT NOT NULL REFERENCES "roles" ("role_id") ON DELETE CASCADE,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);

CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

CREATE INDEX idx_user_roles_user_id_role_id ON user_roles (user_id, role_id);

CREATE INDEX idx_user_roles_created_at ON user_roles (created_at);

CREATE INDEX idx_user_roles_updated_at ON user_roles (updated_at);

CREATE TABLE refresh_tokens (
    refresh_token_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expiration TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);

CREATE INDEX idx_refresh_tokens_expiration ON refresh_tokens (expiration);
