use thiserror::Error;

#[derive(Debug, Error, PartialEq, Eq)]
pub enum GeneratorError {
    #[error("node ID must be between 0 and {maximum}")]
    InvalidNodeId { maximum: u16 },

    #[error("current time is earlier than the custom epoch")]
    TimeBeforeEpoch,

    #[error("system clock moved backwards from {previous_ms} to {current_ms}")]
    ClockMovedBackwards { previous_ms: u64, current_ms: u64 },

    #[error("system clock could not produce a valid Unix timestamp")]
    InvalidSystemTime,
}
