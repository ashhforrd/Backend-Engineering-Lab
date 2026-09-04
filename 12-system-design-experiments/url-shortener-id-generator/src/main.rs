use url_shortener_id_generator::{clock::SystemClock, generator::IdGenerator};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generator = IdGenerator::new(1, SystemClock)?;

    for _ in 0..5 {
        let generated = generator.generate()?;

        println!(
            "short={} numeric={} node={} sequence={}",
            generated.short, generated.numeric, generated.node_id, generated.sequence,
        );
    }

    Ok(())
}
