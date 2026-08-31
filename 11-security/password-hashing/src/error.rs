use thiserror::Error;

#[derive(Debug, Error, PartialEq, Eq)]
pub enum PasswordError {
    #[error("password must contain at least 12 characters")]
    TooShort,

    #[error("password must contain at least one uppercase letter")]
    MissingUppercase,

    #[error("password must contain at least one lowercase letter")]
    MissingLowercase,

    #[error("password must contain at least one digit")]
    MissingDigit,

    #[error("password hashing failed")]
    HashingFailed,

    #[error("stored password hash is invalid")]
    InvalidStoredHash,
}
