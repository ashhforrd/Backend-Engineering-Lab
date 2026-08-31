use argon2::{
    Argon2,
    password_hash::{PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
};
use rand_core::OsRng;

use crate::{error::PasswordError, policy::PasswordPolicy, secret::SecretPassword};

pub struct PasswordService {
    policy: PasswordPolicy,
    argon2: Argon2<'static>,
}

impl PasswordService {
    pub fn new(policy: PasswordPolicy) -> Self {
        Self {
            policy,
            argon2: Argon2::default(),
        }
    }

    pub fn hash(&self, password: &SecretPassword) -> Result<String, PasswordError> {
        self.policy.validate(password.expose())?;

        let salt = SaltString::generate(&mut OsRng);

        self.argon2
            .hash_password(password.expose().as_bytes(), &salt)
            .map(|hash| hash.to_string())
            .map_err(|_| PasswordError::HashingFailed)
    }

    pub fn verify(
        &self,
        password: &SecretPassword,
        stored_hash: &str,
    ) -> Result<bool, PasswordError> {
        let parsed_hash =
            PasswordHash::new(stored_hash).map_err(|_| PasswordError::InvalidStoredHash)?;

        Ok(self
            .argon2
            .verify_password(password.expose().as_bytes(), &parsed_hash)
            .is_ok())
    }
}
