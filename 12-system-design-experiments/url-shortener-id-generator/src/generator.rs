use std::sync::Mutex;

use crate::{
    clock::Clock,
    error::GeneratorError,
};

const NODE_ID_BITS: u8 = 10;
const SEQUENCE_BITS: u8 = 12;

const MAX_NODE_ID: u16 =
    (1_u16 << NODE_ID_BITS) - 1;

struct GeneratorState {
    last_timestamp_ms: u64,
    sequence: u16,
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
    pub fn new(
        node_id: u16,
        clock: C,
    ) -> Result<Self, GeneratorError> {
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
}