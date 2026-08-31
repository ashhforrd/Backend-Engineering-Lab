use password_hashing::{hasher::PasswordService, policy::PasswordPolicy, secret::SecretPassword};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let service = PasswordService::new(PasswordPolicy::default());

    let password = SecretPassword::new("SecurePassword123");

    let stored_hash = service.hash(&password)?;

    let verified = service.verify(&password, &stored_hash)?;

    println!("password verified: {verified}");

    Ok(())
}
