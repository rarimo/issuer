-- +migrate Up

CREATE TABLE committed_states(
    id                    BIGSERIAL    PRIMARY KEY       NOT NULL,
    tx_id                 CHAR(66),
    created_at            TIMESTAMP    WITHOUT TIME ZONE NOT NULL,
    block_timestamp       BIGINT,
    block_number          BIGINT,
    is_genesis            BOOLEAN,
    roots_tree_root       BYTEA                          NOT NULL,
    claims_tree_root      BYTEA                          NOT NULL,
    revocations_tree_root BYTEA                          NOT NULL,
    status                VARCHAR(256)                   NOT NULL,
    message               VARCHAR(256)
);

CREATE TABLE claims(
   id                CHAR(36)     PRIMARY KEY NOT NULL,
   schema_type       VARCHAR(256)             NOT NULL,
   revoked           BOOLEAN,
   data              BYTEA,
   core_claim        BYTEA                    NOT NULL,
   signature_proof   BYTEA,
   user_id           VARCHAR(42)              NOT NULL
);

CREATE TABLE claims_offers(
    id          CHAR(36)     PRIMARY KEY       NOT NULL,
    from_id     VARCHAR(42)                    NOT NULL,
    to_id       CHAR(42)                       NOT NULL,
    created_at  TIMESTAMP    WITHOUT TIME ZONE NOT NULL,
    claim_id    VARCHAR(256)                   NOT NULL,
    is_received BOOLEAN
);

CREATE TABLE roots_tree(
    key   BYTEA PRIMARY KEY NOT NULL,
    value BYTEA             NOT NULL
);

CREATE TABLE revocation_tree(
    key   BYTEA PRIMARY KEY NOT NULL,
    value BYTEA             NOT NULL
);

CREATE TABLE claims_tree(
    key   BYTEA PRIMARY KEY NOT NULL,
    value BYTEA             NOT NULL
);

-- +migrate Down

DROP TABLE committed_states;
DROP TABLE claims_offers;
DROP TABLE claims;
DROP TABLE roots_tree;
DROP TABLE revocation_tree;
DROP TABLE claims_tree;
