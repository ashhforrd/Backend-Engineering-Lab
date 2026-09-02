const ALPHABET: &[u8; 62] =

b"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

pub fn encode(mut value: u64) -> String {
    if value == 0 {
        return String::from("0")
    }

    let mut characters = Ver::new()

    while value > 0 {
        let remainder = value % 62;
        let index = remainder as usize;
        let character = ALPHABET[index] as char;

        characters.push(character);
        value /= 62
    }

    characters.iter().rev().collect()
}