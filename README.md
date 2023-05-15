# issuer

## Description

This is an Issuer service that follows an [Iden3 protocol](https://docs.iden3.io/).
It provides an ability to create an Identity and issue different claims.
You can think of a claim as a statement: something an Issuer says about another subject.
For example, when you apply for a job, the employer may ask you to provide a reference from your previous employer.

1) [Identity generation](#identity-generation)
2) [Claim issuance](#claim-issuance)

### Identity generation

Roughly says the **Identity** itself is a **public/private key pair**. It is generated using
[BabyJubJub elliptic curve](https://eips.ethereum.org/EIPS/eip-2494). But you don't operate with your public
key directly. Instead, you use the [Auth claim](https://docs.iden3.io/protocol/bjjkey/). It is a claim that contains
information about your public key (the _X_, _Y_ in the index fields of the public key) and is signed by the issuer which
is Identity itself. So you can think of the Auth claim as proof that you own the private key.

Every identity has a unique **identifier**. Identifier structure: `Base58 [ type | genesis_state | checksum ]` 
- type: 2 bytes specifying the type 
- genesis_state: first 27 bytes from the identity state (using the **Genesis Claim Merkle tree**) 
- checksum: addition (with overflow) of all the ID bytes Little Endian 16 bits ( `[ type | genesis_state ]`).

**Identity state** is a hash of the three merkle roots: **Claims merkle tree**, **Revocations merkle tree**, and 
**Roots merkle tree**. Genesis **Claims merkle tree** contains only the **Auth claim**. 

Identifier example: `115n3Lx26aHrLfcAWtYnoamo5FDAzWb9PnjMndwtRe`

Every leaf of the **Claims merkle tree** contains the claim's index fields hash as an ID and the claim's value fields hash as a value.

<details>
<summary>General Claim structure</summary>

```
h_i = H(i_0, i_1, i_2, i_3)
h_v = H(v_0, v_1, v_2, v_3)
h_t = H(h_i, h_v)

Index:
 i_0: [ 128 bits ] claim schema
      [ 32 bits ] header flags
          [3] Subject:
            000: A.1 Self
            001: invalid
            010: A.2.i OtherIden Index
            011: A.2.v OtherIden Value
            100: B.i Object Index
            101: B.v Object Value
          [1] Expiration: bool
          [1] Updatable: bool
          [27] 0
      [ 32 bits ] version (optional?)
      [ 61 bits ] 0 - reserved for future use
 i_1: [ 248 bits] identity (case b) (optional)
      [  5 bits ] 0
 i_2: [ 253 bits] 0
 i_3: [ 253 bits] 0
Value:
 v_0: [ 64 bits ]  revocation nonce
         [ 64 bits ]  expiration date (optional)
         [ 125 bits] 0 - reserved
 v_1: [ 248 bits] identity (case c) (optional)
        [  5 bits ] 0
 v_2: [ 253 bits] 0
 v_3: [ 253 bits] 0
```

</details>

<details>
<summary>Auth claim structure</summary>

```
Index:
 i_0: [ 128 bits] 269270088098491255471307608775043319525 // auth schema (big integer from ca938857241db9451ea329256b9c06e5)
      [ 32 bits ] 00010000000000000000 // header flags: first 000 - self claim 1 - expiration is set. 
      [ 32 bits ] 0
      [ 61 bits ] 0 
 i_1: [ 253 bits] 0
 i_2: [ 253 bits] 15730379921066174438220083697399546667862601297001890929936158339406931652649 // x part of BJJ pubkey
 i_3: [ 253 bits] 5635420193976628435572861747946801377895543276711153351053385881432935772762  // y part of BJJ pubkey
Value:
 v_0: [ 64 bits ] 2484496687 // revocation nonce
      [ 64 bits ] 1679670808 // expiration timestamp
      [ 125 bits] 0
 v_1: [ 253 bits] 0
 v_2: [ 253 bits] 0
 v_3: [ 253 bits] 0
```

</details>

For describing the claim we use [JSON-ld schemas](https://json-ld.org/). It is provided an easy way to tie the claim, 
it's fields to descriptions.

<details>
<summary>Auth claim schema</summary>

```json
{
  "@context": [{
    "@version": 1.1,
    "@protected": true,
    "id": "@id",
    "type": "@type",
    "AuthBJJCredential": {
      "@id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld#AuthBJJCredential",
      "@context": {
        "@version": 1.1,
        "@protected": true,
        "id": "@id",
        "type": "@type",
        "auth-vocab": "https://github.com/iden3/claim-schema-vocab/blob/main/credentials/auth.md#",
        "serialization": "https://github.com/iden3/claim-schema-vocab/blob/main/credentials/serialization.md#",
        "x": {
          "@id": "auth-vocab:x",
          "@type": "serialization:IndexDataSlotA"
        },
        "y": {
          "@id": "auth-vocab:y",
          "@type": "serialization:IndexDataSlotB"
        }
      }
    }
  }]
}
```    

</details>

### Claim issuance

For the claim issuance included the following steps:
1) KYC service verify the data provided by the user and if it is correct - send the request to the issuer.
2) Issuer create a claim by schema with user information that already verified by the KYC service.
3) User retrieves the claim.

Example of the claim issuance flow:
``` mermaid

sequenceDiagram
    participant User
    participant KYC Service
    participant Issuer

    activate User
    
    User->>KYC Service: Provide a claim data with <br> ZKP of the Identity owning
    activate KYC Service
    KYC Service->>KYC Service: Verify the claim data
    KYC Service->>Issuer: Send request for <br> claim issuance
    deactivate KYC Service
    
    activate Issuer
    Issuer->>Issuer: Create claim by schema <br> with the user information
    Issuer->>Issuer: Add claim to the <br> Claims merkle Tree
    Issuer->>Issuer: Publish on-chain <br> the updated state
    deactivate Issuer
    
    User->>Issuer: Request the claim with <br> ZKP of the Identity owning
    activate Issuer
    Issuer->>Issuer: Create MTP and Sign proof <br> of the claim storing in the <br> Issuer's Claims merkle Tree
    Issuer->>User: Send the claim with <br> ZKP of the Issuer's Identity <br> owning, MTP and Sign proof
    
    deactivate Issuer
   
    deactivate User
    
```

### MTP (Merkle Tree Proof)

Where MTP (Merkle Tree Proof) is a proof of the claim storing in the Issuer's Claims merkle Tree. MTP contains from
the all siblings of the claim leaf in the tree that provide an ability with hashing recalculate the root hash of the tree.
And in the end compare it with the root hash that stored on-chain. If the hashes are equal - the claim is valid.

<details>
<summary>MTP structure</summary>

```json
{
    "@type": "Iden3SparseMerkleProof",
    "issuer_data": {
    "id": "115n3Lx26aHrLfcAWtYnoamo5FDAzWb9PnjMndwtRe",
    "state": {
        "block_number": 4700047,
        "block_timestamp": 1675445798,
        "claims_tree_root": "6db7cd72af198b3d87a96cc226e6252c38168d41c130449760c33f1e65ce721d",
        "revocation_tree_root": "0000000000000000000000000000000000000000000000000000000000000000",
        "root_of_roots": "aea53806155ce649d04903da174b22a2eeff758f65499407427dfca88ef71612",
        "tx_id": "0xeaceb1c2e689a7cc1730d8b42215c0250f6945d868b7a9f741bf273bfeadef59",
        "value": "1034dbf7e4a6808220f0d1b2872923082bc87702b46d7e747d15fe40431fba11"
    }
    },
    "mtp": {
        "existence": true,
        "siblings": [
            "0",
            "0",
            "0",
            "15225701537283030212784216827065714165069625960211407334852750198781065225133"
        ]
    }
}
```

</details>

### Sign proof

Sign proof consists of the Issuer's signature of the claim and the proof that Auth claim is stored in the 
Issuer's Claims Merkle Tree.

<details>
<summary>Sign proof structure</summary>

```json
{
    "@type": "BJJSignature2021",
    "issuer_data": {
        "auth_claim": [
            "304427537360709784173770334266246861770",
            "0",
            "856628192321641508866307593746150915558488491935060423534997169059498987478",
            "9494420313137459273421780860424623914451830428162908508853978141159438399464",
            "384513778",
            "0",
            "0",
            "0"
        ],
        "id": "115n3Lx26aHrLfcAWtYnoamo5FDAzWb9PnjMndwtRe",
        "mtp": {
            "existence": true,
            "siblings": [
                "0",
                "0",
                "0",
                "164064612995325705517108283320576642342848431716673310363948244017845964996"
            ]
        },
        "revocation_status": "https://8c44-193-193-222-99.eu.ngrok.io/integrations/issuer/v1/claims/revocations/check/2189031102",
        "state": {
            "claims_tree_root": "6db7cd72af198b3d87a96cc226e6252c38168d41c130449760c33f1e65ce721d",
            "value": "1034dbf7e4a6808220f0d1b2872923082bc87702b46d7e747d15fe40431fba11"
        }
    },
    "signature": "4facd0b7181a903f9558cda9ca9145c9ef8455f663330479ab70869b1bcba70ba4d510ae4d0dffe6a006de81467bbe6030f7e0a338881126e0298a070fb29005"
}
```

</details>

## Install

  ```
  git clone <repo_url>/issuer
  cd issuer
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main migrate up
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.


## Running from docker 
  
Make sure that docker installed.

use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t gitlab.com/q-dev/q-id/issuer .
  docker run -e KV_VIPER_FILE=/config.yaml gitlab.com/q-dev/q-id/issuer
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command


### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).

