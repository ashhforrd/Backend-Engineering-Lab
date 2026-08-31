use crate::error::PasswordError;

pub struct PasswordPolicy {
    minimum_length: usize,
}

impl PasswordPolicy {
    pub fn new(minimum_length: usize) -> Self {
        Self { minimum_length }
    }

    pub fn validate(&self, password: &str) -> Result<(), PasswordError> {
        if password.chars().count() < self.minimum_length {
            return Err(PasswordError::TooShort);
        }

        if !password.chars().any(char::is_uppercase) {
            return Err(PasswordError::MissingUppercase);
        }

        if !password.chars().any(char::is_lowercase) {
            return Err(PasswordError::MissingLowercase);
        }

        if !password.chars().any(|character| character.is_ascii_digit()) {
            return Err(PasswordError::MissingDigit);
        }

        Ok(())
    }
}

impl Default for PasswordPolicy {
    fn default() -> Self {
        Self { minimum_length: 12 }
    }
}
