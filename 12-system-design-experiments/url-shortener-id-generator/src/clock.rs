use std::time::{SystemTime, UNIX_EPOCH};

use crate::error::GeneratorError;

pub trait Clock {
    fn now_millis(&self) -> Result<u64, GeneratorError>;
}

pub struct SystemClock;

impl Clock for SystemClock {
    fn now_millis(&self) -> Result<u64, GeneratorError> {
        let duration = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map_err(|_| GeneratorError::InvalidSystemTime)?;

        u64::try_from(duration.as_millis()).map_err(|_| GeneratorError::InvalidSystemTime)
    }
}
