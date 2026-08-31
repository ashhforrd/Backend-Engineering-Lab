use password_hashing::{
    error::PasswordError, hasher::PasswordService, policy::PasswordPolicy, secret::SecretPassword,
};

fn service() -> PasswordService {
    PasswordService::new(PasswordPolicy::default())
}

#[test]
fn hashes_and_verifies_the_correct_password() {
    let service = service();
    let password = SecretPassword::new("SecurePassword123");

    let hash = service.hash(&password).unwrap();

    assert!(service.verify(&password, &hash).unwrap());
}

#[test]
fn rejects_an_incorrect_password() {
    let service = service();
    let correct_password = SecretPassword::new("SecurePassword123");
    let incorrect_password = SecretPassword::new("DifferentPassword456");

    let hash = service.hash(&correct_password).unwrap();

    assert!(!service.verify(&incorrect_password, &hash).unwrap());
}

#[test]
fn produces_different_hashes_for_the_same_password() {
    let service = service();
    let password = SecretPassword::new("SecurePassword123");

    let first_hash = service.hash(&password).unwrap();
    let second_hash = service.hash(&password).unwrap();

    assert_ne!(first_hash, second_hash);
}

#[test]
fn rejects_a_password_that_violates_the_policy() {
    let service = service();
    let password = SecretPassword::new("short1A");

    let result = service.hash(&password);

    assert_eq!(result, Err(PasswordError::TooShort));
}

#[test]
fn rejects_an_invalid_stored_hash() {
    let service = service();
    let password = SecretPassword::new("SecurePassword123");

    let result = service.verify(&password, "not-a-valid-hash");

    assert_eq!(result, Err(PasswordError::InvalidStoredHash));
}
