use std::{sync::Mutex, thread};

use crate::{base62, clock::Clock, error::GeneratorError};

const CUSTOM_EPOCH_MS: u64 = 1_767_225_600_000;

const TIMESTAMP_BITS: u8 = 41;
const NODE_ID_BITS: u8 = 10;
const SEQUENCE_BITS: u8 = 12;

const MAX_TIMESTAMP: u64 = (1_u64 << TIMESTAMP_BITS) - 1;

const MAX_NODE_ID: u16 = (1_u16 << NODE_ID_BITS) - 1;

const MAX_SEQUENCE: u16 = (1_u16 << SEQUENCE_BITS) - 1;

struct GeneratorState {
    last_timestamp_ms: u64,
    sequence: u16,
}

#[derive(Debug, PartialEq, Eq)]
pub struct GeneratedId {
    pub numeric: u64,
    pub short: String,
    pub timestamp_ms: u64,
    pub node_id: u16,
    pub sequence: u16,
}

pub struct IdGenerator<C>
where
    C: Clock,
{
    node_id: u16,
    clock: C,
    state: Mutex<GeneratorState>,
}

impl<C> IdGenerator<C>
where
    C: Clock,
{
    pub fn new(node_id: u16, clock: C) -> Result<Self, GeneratorError> {
        if node_id > MAX_NODE_ID {
            return Err(GeneratorError::InvalidNodeId {
                maximum: MAX_NODE_ID,
            });
        }

        Ok(Self {
            node_id,
            clock,
            state: Mutex::new(GeneratorState {
                last_timestamp_ms: 0,
                sequence: 0,
            }),
        })
    }

    pub fn generate(&self) -> Result<GeneratedId, GeneratorError> {
        loop {
            let timestamp_ms = self.clock.now_millis()?;

            if timestamp_ms < CUSTOM_EPOCH_MS {
                return Err(GeneratorError::TimeBeforeEpoch);
            }

            let elapsed_ms = timestamp_ms - CUSTOM_EPOCH_MS;

            if elapsed_ms > MAX_TIMESTAMP {
                return Err(GeneratorError::TimestampExhausted);
            }

            let mut state = self
                .state
                .lock()
                .map_err(|_| GeneratorError::StateUnavailable)?;

            if timestamp_ms < state.last_timestamp_ms {
                return Err(GeneratorError::ClockMovedBackwards {
                    previous_ms: state.last_timestamp_ms,
                    current_ms: timestamp_ms,
                });
            }

            let sequence = if timestamp_ms == state.last_timestamp_ms {
                if state.sequence == MAX_SEQUENCE {
                    drop(state);
                    thread::yield_now();
                    continue;
                }

                state.sequence += 1;
                state.sequence
            } else {
                state.last_timestamp_ms = timestamp_ms;
                state.sequence = 0;
                state.sequence
            };

            let numeric = (elapsed_ms << (NODE_ID_BITS + SEQUENCE_BITS))
                | (u64::from(self.node_id) << SEQUENCE_BITS)
                | u64::from(sequence);

            return Ok(GeneratedId {
                numeric,
                short: base62::encode(numeric),
                timestamp_ms,
                node_id: self.node_id,
                sequence,
            });
        }
    }
}
