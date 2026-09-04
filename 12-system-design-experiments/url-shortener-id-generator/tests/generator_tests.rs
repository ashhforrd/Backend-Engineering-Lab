use std::{
    collections::HashSet,
    sync::{
        Arc,
        atomic::{AtomicU64, Ordering},
    },
    thread,
};

use url_shortener_id_generator::{
    base62, clock::Clock, error::GeneratorError, generator::IdGenerator,
};

const EPOCH_MS: u64 = 1_767_225_600_000;

#[derive(Clone)]
struct TestClock {
    timestamp_ms: Arc<AtomicU64>,
}

impl TestClock {
    fn new(timestamp_ms: u64) -> Self {
        Self {
            timestamp_ms: Arc::new(AtomicU64::new(timestamp_ms)),
        }
    }

    fn set(&self, timestamp_ms: u64) {
        self.timestamp_ms.store(timestamp_ms, Ordering::SeqCst);
    }
}

impl Clock for TestClock {
    fn now_millis(&self) -> Result<u64, GeneratorError> {
        Ok(self.timestamp_ms.load(Ordering::SeqCst))
    }
}

#[test]
fn encodes_numbers_as_base62() {
    assert_eq!(base62::encode(0), "0");
    assert_eq!(base62::encode(61), "Z");
    assert_eq!(base62::encode(62), "10");
    assert_eq!(base62::encode(125), "21");
}

#[test]
fn rejects_a_node_id_that_does_not_fit_in_ten_bits() {
    let result = IdGenerator::new(1024, TestClock::new(EPOCH_MS));

    assert_eq!(
        result.err(),
        Some(GeneratorError::InvalidNodeId { maximum: 1023 }),
    );
}

#[test]
fn increments_sequence_within_the_same_millisecond() {
    let generator = IdGenerator::new(7, TestClock::new(EPOCH_MS + 1)).unwrap();

    let first = generator.generate().unwrap();
    let second = generator.generate().unwrap();

    assert_eq!(first.sequence, 0);
    assert_eq!(second.sequence, 1);
    assert_ne!(first.numeric, second.numeric);
    assert_ne!(first.short, second.short);
}

#[test]
fn resets_sequence_when_time_moves_forward() {
    let clock = TestClock::new(EPOCH_MS + 1);
    let generator = IdGenerator::new(7, clock.clone()).unwrap();

    generator.generate().unwrap();
    generator.generate().unwrap();
    clock.set(EPOCH_MS + 2);

    let generated = generator.generate().unwrap();

    assert_eq!(generated.sequence, 0);
}

#[test]
fn rejects_a_clock_that_moves_backwards() {
    let clock = TestClock::new(EPOCH_MS + 2);
    let generator = IdGenerator::new(7, clock.clone()).unwrap();

    generator.generate().unwrap();
    clock.set(EPOCH_MS + 1);

    let result = generator.generate();

    assert_eq!(
        result,
        Err(GeneratorError::ClockMovedBackwards {
            previous_ms: EPOCH_MS + 2,
            current_ms: EPOCH_MS + 1,
        }),
    );
}

#[test]
fn generates_unique_ids_across_concurrent_threads() {
    let generator = Arc::new(IdGenerator::new(7, TestClock::new(EPOCH_MS + 1)).unwrap());

    let handles: Vec<_> = (0..8)
        .map(|_| {
            let generator = Arc::clone(&generator);

            thread::spawn(move || {
                (0..100)
                    .map(|_| generator.generate().unwrap().numeric)
                    .collect::<Vec<_>>()
            })
        })
        .collect();

    let ids: Vec<u64> = handles
        .into_iter()
        .flat_map(|handle| handle.join().unwrap())
        .collect();

    let unique_ids: HashSet<u64> = ids.iter().copied().collect();

    assert_eq!(ids.len(), 800);
    assert_eq!(unique_ids.len(), 800);
}
