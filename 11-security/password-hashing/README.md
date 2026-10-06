# Password Hashing

## Problem

Passwords must not be stored as plaintext or as fast general-purpose hashes. A
database leak should not immediately reveal reusable credentials.

## Design

The service validates password policy, generates a random salt, hashes with
Argon2id, and stores the complete PHC string. Verification parses that string
and compares a submitted password without recovering the original value.

## Implementation

- Rust and the RustCrypto `argon2` crate
- Argon2id with encoded algorithm parameters
- Operating-system randomness for unique salts
- Dedicated password policy and typed errors
- Secret wrapper that zeroizes owned password memory on drop
- Tests for verification, invalid hashes, policy, and salt uniqueness

## Failure Cases

- Weak passwords are rejected before hashing.
- Malformed stored hashes produce a typed error.
- Incorrect passwords return `false` without revealing which part differed.
- Zeroization cannot erase copies created outside the secret wrapper.
- Hash parameters must be reviewed as hardware and security guidance evolve.

## What I Learned

- Password hashing should be intentionally expensive and memory-hard.
- A salt prevents identical passwords from producing identical stored hashes.
- The PHC string carries the salt and algorithm parameters with the hash.
- Memory safety helps, but secret lifetime and accidental copies still matter.

## Running

```bash
cargo run
cargo test
```
